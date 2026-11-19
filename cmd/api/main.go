package main

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/riyanamanda/helpdesk-backend/internal/platform/presence"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rustfs"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/profile"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
	"github.com/riyanamanda/helpdesk-backend/internal/simgos"
	"github.com/riyanamanda/helpdesk-backend/internal/ticket"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()
	cfg := config.Load()

	slog.Info("config loaded",
		"addr", net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		"simgos_enabled", cfg.IhsDatabase.Host != "",
		"bpjs_enabled", cfg.BPJS.VClaimURL != "",
		"antrol_enabled", cfg.Antrol.Domain != "",
		"google_login_enabled", cfg.Auth.FirebaseProjectID != "",
	)
	if names := cfg.InsecureDefaults(); len(names) > 0 {
		slog.Warn("insecure default config in use", "env", names)
	}

	// Database
	slog.Info("connecting to database")
	db, err := database.NewPostgres(cfg.Database.ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	slog.Info("running migrations")
	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("run migration: %w", err)
	}

	txManager := database.NewTxManager(db)

	// SIMGOS Database
	var simgosDB *sqlx.DB

	if cfg.IhsDatabase.Host != "" {
		slog.Info("connecting to simgos database")
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
		return fmt.Errorf("init bucket: %w", err)
	}

	storageService := storage.NewRustFSStorage(rustfsClient, cfg.Storage.Bucket)

	// RabbitMQ
	slog.Info("connecting to rabbitmq")
	rabbitmqClient, err := rabbitmq.Connect(cfg.RabbitMQ.RabbitMQConnString())
	if err != nil {
		return fmt.Errorf("connect rabbitmq: %w", err)
	}
	defer rabbitmqClient.Close()

	// Redis
	slog.Info("connecting to redis")
	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		return err
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

	presenceTracker := presence.NewTracker(redisClient, userRepo)

	// Services
	categoryService := category.NewService(categoryRepo, cacheStore)
	divisionService := division.NewService(divisionRepo, cacheStore)
	userService := user.NewService(userRepo, outboxRepo, txManager, cfg.Storage, cacheStore, presenceTracker)
	ticketService := ticket.NewService(ticketRepo, outboxRepo, txManager, storageService, cfg.Storage, cacheStore, categoryService, divisionService, userService)
	feedbackService := feedback.NewService(feedbackRepo)
	profileService := profile.NewService(profileRepo, storageService, cfg.Storage, cfg.Auth, cacheStore)
	rbacService := rbac.NewService(rbacRepo, cacheStore)
	permissionService := rbac.NewPermissionService(rbacRepo, cacheStore)
	dashboardService := dashboard.NewService(dashboardRepo, cacheStore)

	ihsService := ihs.NewService(ihsRepo)
	antrianService := antrian.NewService(antrianRepo, simgosClient)
	authService := auth.NewService(userRepo, cfg.Auth, cfg.Storage, cfg.App, cacheStore, permissionService, *rabbitmqClient, presenceTracker)
	bpjsService := bpjs.NewService(bpjsClient)

	// HTTP
	e := echo.New()
	e.Validator = validation.New()
	e.HTTPErrorHandler = response.HTTPErrorHandler

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
	authMiddleware := middleware.AuthMiddleware(cfg.Auth, redisClient, permissionService, presenceTracker)
	protected.Use(authMiddleware)
	auth.Register(api, authService, authMiddleware)
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

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("start server: %w", err)
	case <-quit:
	}

	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	slog.Info("server stopped")

	return nil
}
