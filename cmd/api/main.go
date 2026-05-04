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
	"github.com/labstack/echo/v5/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/config"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// health check route
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"name":   cfg.AppName,
		})
	})

	// depencencies

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.AppHost, cfg.AppPort),
		Handler: e,
	}

	// start server in goroutine
	go func() {
		slog.Info("server starting", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error to start server", "host", server.Addr)
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
