package main

import (
	"log/slog"
	"os"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/database"
	"github.com/riyanamanda/helpdesk-backend/internal/seed"
)

func main() {
	slog.Info("starting database seeding")

	cfg := config.Load()
	db := database.NewPostgres(cfg.Database.ConnString())
	defer db.Close()

	if err := seed.Run(db); err != nil {
		slog.Error("database seeding failed", "error", err)
		os.Exit(1)
	}

	slog.Info("database seeding completed")
}
