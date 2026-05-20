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
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/database"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/ticket"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.Validator = validation.New()

	middleware.Register(e)

	// health check route
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"name":   cfg.App.Name,
		})
	})

	// dependencies
	// postgres init
	db := database.NewPostgres(cfg.Database.ConnString())
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	// minio storage init
	minioClient, err := storage.NewMinioClient(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey, cfg.Storage.UseSSL)
	if err != nil {
		slog.Error("load storage failed", "error", err)
		os.Exit(1)
	}
	storageService := storage.NewMinioStorage(minioClient, cfg.Storage.Bucket, cfg.Storage.PublicURL)

	// api root
	api := e.Group("/api/v1")

	// public routes
	auth.Register(api, db, cfg.Auth, cfg.Storage)

	// protected route
	protected := api.Group("")
	protected.Use(
		middleware.AuthMiddleware(cfg.Auth),
	)

	// protected modules
	category.Register(protected, db)
	division.Register(protected, db)
	user.Register(protected, db, storageService, cfg.Storage)
	ticket.Register(protected, db, storageService, cfg.Storage)

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: e,
	}

	// start server in goroutine
	go func() {
		slog.Info("server starting", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "addr", server.Addr, "error", err)
		}
	}()

	// channel to capture signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // blocking
	slog.Info("shutting down server...")

	// wait 10sec for finish any request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	slog.Info("server exited properly")
}
