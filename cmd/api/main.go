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

	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/auth"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/minio"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/profile"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/ticket"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	e := echo.New()
	e.Validator = validation.New()

	middleware.Register(e, cfg.App)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"name":   cfg.App.Name,
		})
	})

	db := database.NewPostgres(cfg.Database.ConnString())
	defer db.Close()

	minioClient, err := minio.NewMinioClient(
		cfg.Storage.Endpoint,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.UseSSL,
	)
	if err != nil {
		slog.Error("load storage failed", "error", err)
		os.Exit(1)
	}

	if err := minio.InitBucket(ctx, minioClient, cfg.Storage.Bucket); err != nil {
		slog.Error("failed to initialize storage bucket", "error", err)
		os.Exit(1)
	}

	storageService := storage.NewMinioStorage(
		minioClient,
		cfg.Storage.Bucket,
	)

	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		slog.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}

	cacheStore := cache.NewRedisCache(redisClient)

	userRepo := user.NewUserRepository(db)
	api := e.Group("/api/v1")

	auth.Register(api, userRepo, cfg.Auth, cfg.Storage, redisClient)

	protected := api.Group("")
	protected.Use(
		middleware.AuthMiddleware(cfg.Auth, redisClient),
	)

	category.Register(protected, db, cacheStore)
	division.Register(protected, db, cacheStore)
	user.Register(protected, userRepo, cfg.Storage, cacheStore)
	ticket.Register(protected, db, storageService, cfg.Storage, cacheStore)
	dashboard.Register(protected, db, cacheStore)
	profile.Register(protected, db, storageService, cfg.Storage, cfg.Auth)

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
