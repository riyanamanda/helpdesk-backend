package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/antrian"
	"github.com/riyanamanda/helpdesk-backend/internal/auth"
	"github.com/riyanamanda/helpdesk-backend/internal/bpjs"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/feedback"
	"github.com/riyanamanda/helpdesk-backend/internal/ihs"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rustfs"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/profile"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
	"github.com/riyanamanda/helpdesk-backend/internal/simgos"
	"github.com/riyanamanda/helpdesk-backend/internal/ticket"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()
	cfg := config.Load()

	// Database
	slog.Info("connecting to database")
	db := database.NewPostgres(cfg.Database.ConnString())
	defer db.Close()

	slog.Info("running migrations")
	if err := database.RunMigrations(db); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	txManager := database.NewTxManager(db)

	// SIMGOS Database
	var simgosDB *sqlx.DB

	if cfg.IhsDatabase.Host != "" {
		slog.Info("connecting to simgos database")

		var err error
		simgosDB, err = database.NewMySql(cfg.IhsDatabase.MySqlConnString())
		if err != nil {
			slog.Warn("simgos database unavailable, simgos routes disabled", "error", err)
			simgosDB = nil
		}
	} else {
		slog.Warn("simgos database not configured, simgos routes disabled")
	}

	if simgosDB != nil {
		defer simgosDB.Close()
	}

	// RustFS
	slog.Info("connecting to object storage")

	rustfsClient := rustfs.NewClient(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey)

	if err := rustfs.InitBucket(ctx, rustfsClient, cfg.Storage.Bucket); err != nil {
		slog.Error("rustfs bucket initialization failed", "error", err)
		os.Exit(1)
	}

	storageService := storage.NewRustFSStorage(rustfsClient, cfg.Storage.Bucket)

	// RabbitMQ
	slog.Info("connecting to RabbitMQ")

	rabbitmqClient, err := rabbitmq.Connect(cfg.RabbitMQ.RabbitMQConnString())
	if err != nil {
		slog.Error("RabbitMQ connection failed", "error", err)
		os.Exit(1)
	}
	defer rabbitmqClient.Close()

	// Redis
	slog.Info("connecting to Redis")

	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		slog.Error("Redis connection failed", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	cacheStore := cache.NewRedisCache(redisClient)

	// simgos
	simgosClient := simgos.NewAntrolClient(cfg.Antrol.Domain, cfg.Antrol.Username, cfg.Antrol.Password)

	// BPJS
	bpjsClient := bpjs.NewClient(cfg.BPJS.VClaimURL, cfg.BPJS.ConsID, cfg.BPJS.ConsSecret, cfg.BPJS.VClaimKey)

	// Repositories
	categoryRepo := category.NewRepository(db)
	divisionRepo := division.NewRepository(db)
	userRepo := user.NewRepository(db)
	ticketRepo := ticket.NewRepository(db)
	feedbackRepo := feedback.NewRepository(db)
	profileRepo := profile.NewRepository(db)
	rbacRepo := rbac.NewRepository(db)
	dashboardRepo := dashboard.NewRepository(db)
	outboxRepo := outbox.NewRepository(db)

	ihsRepo := ihs.NewRepository(simgosDB)
	antrianRepo := antrian.NewRepository(simgosDB)

	// Services
	categoryService := category.NewService(categoryRepo, cacheStore)
	divisionService := division.NewService(divisionRepo, cacheStore)
	userService := user.NewService(userRepo, outboxRepo, txManager, cfg.Storage, cacheStore)
	ticketService := ticket.NewService(ticketRepo, outboxRepo, txManager, storageService, cfg.Storage, cacheStore, categoryService, divisionService, userService)
	feedbackService := feedback.NewService(feedbackRepo)
	profileService := profile.NewService(profileRepo, storageService, cfg.Storage, cfg.Auth)
	rbacService := rbac.NewService(rbacRepo, cacheStore)
	permissionService := rbac.NewPermissionService(rbacRepo, cacheStore)
	dashboardService := dashboard.NewService(dashboardRepo, cacheStore)

	ihsService := ihs.NewService(ihsRepo)
	antrianService := antrian.NewService(antrianRepo, simgosClient)
	authService := auth.NewService(userRepo, cfg.Auth, cfg.Storage, cfg.App, cacheStore, permissionService, *rabbitmqClient)
	bpjsService := bpjs.NewService(bpjsClient)

	// HTTP
	e := echo.New()
	e.Validator = validation.New()

	middleware.Register(e, cfg.App)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"name":   cfg.App.Name,
		})
	})

	e.Match([]string{http.MethodGet, http.MethodHead}, "/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API
	api := e.Group("/api/v1")
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.Auth, redisClient, permissionService))
	auth.Register(api, authService, middleware.AuthMiddleware(cfg.Auth, redisClient, permissionService))
	category.Register(protected, categoryService)
	division.Register(protected, divisionService)
	user.Register(protected, userService)
	ticket.Register(protected, ticketService)
	feedback.Register(protected, feedbackService)
	profile.Register(protected, profileService)
	rbac.Register(protected, rbacService)
	dashboard.Register(protected, dashboardService)
	bpjs.Register(protected, bpjsService)

	ihs.Register(protected, ihsService)
	if simgosDB != nil {
		antrian.Register(protected, antrianService)
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: e,
	}

	go func() {
		slog.Info("server starting", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "addr", server.Addr, "error", err)
		}
	}()

	// gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	slog.Info("server exited properly")
}
