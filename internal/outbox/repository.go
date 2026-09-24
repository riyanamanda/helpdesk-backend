package outbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type Repository interface {
	Create(ctx context.Context, tx database.Tx, event OutboxEvent) error
	GetUnpublished(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
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
		INSERT INTO outbox_events (event_type, aggregate_id, payload)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`

	_, err := tx.ExecContext(ctx, query, event.EventType, event.AggregateID, event.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetUnpublished(ctx context.Context, limit int) ([]OutboxEvent, error) {
	var events []OutboxEvent

	const query = `
		SELECT
			id,
			event_type,
			aggregate_id,
			payload,
			created_at,
			published_at
		FROM outbox_events
		WHERE published_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
	`

	if err := r.db.SelectContext(ctx, &events, query, limit); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *repository) MarkPublished(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE outbox_events
		SET published_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
