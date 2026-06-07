package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
)

func main() {
	cfg := config.Load()

	cleanup, err := bootstrap(cfg)
	if err != nil {
		slog.Error("failed to bootstrap worker", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down worker...")
}
