package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/seed"
)

func main() {
	if err := run(); err != nil {
		slog.Error("seed stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("starting database seeding")

	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database.ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := seed.Run(db); err != nil {
		return fmt.Errorf("run seed: %w", err)

	}

	slog.Info("database seeding completed")

	return nil
}
