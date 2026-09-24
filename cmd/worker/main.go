package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	db := database.NewPostgres(cfg.Database.ConnString())
	defer db.Close()

	outboxRepo := outbox.NewRepository(db)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			events, err := outboxRepo.GetUnpublished(ctx, 10)
			if err != nil {
				slog.Error("get unpublished events", "error", err)
				continue
			}

			for _, event := range events {
				slog.Info("outbox_events", "id", event.ID, "type", event.EventType, "aggregate_id", event.AggregateID)
			}
		}
	}
}
