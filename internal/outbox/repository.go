package outbox

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type Repository interface {
	Create(ctx context.Context, tx database.Tx, event OutboxEvent) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, tx database.Tx, event OutboxEvent) error {
	const query = `
		INSERT INTO outbox_events (event_type, agregate_id, payload)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`

	_, err := tx.ExecContext(ctx, query, event.EventType, event.AggregateID, event.Payload)
	if err != nil {
		return err
	}

	return nil
}
