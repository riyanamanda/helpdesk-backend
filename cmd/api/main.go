package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	server, cleanup, err := bootstrap(ctx, cfg)
	if err != nil {
		slog.Error("failed to bootstrap application", "error", err)
		os.Exit(1)
	}
	defer cleanup()

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
