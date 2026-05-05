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
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/config"
	"github.com/riyanamanda/helpdesk-backend/internal/database"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		HandleError:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			level := slog.LevelInfo
			if v.Status >= 500 {
				level = slog.LevelError
			}
			slog.Log(c.Request().Context(), level, "http request",
				"request_id", v.RequestID,
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"error", v.Error,
			)

			return nil
		},
	}))
	e.Use(middleware.Recover())

	// health check route
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"name":   cfg.AppName,
		})
	})

	// depencencies
	db := database.NewPostgres(cfg.DBConnString())
	defer db.Close()

	// routes
	api := e.Group("/api/v1")
	category.Register(api, db)

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.AppHost, cfg.AppPort),
		Handler: e,
	}

	// start server in goroutine
	go func() {
		slog.Info("server starting", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error to start server", "addr", server.Addr, "error", err)
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
