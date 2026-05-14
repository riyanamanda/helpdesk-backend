package ticket

import (
	"context"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name TicketRepository
type TicketRepository interface {
	GetAll(ctx context.Context, params GetTicketParams) ([]Ticket, int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetTicketParams) ([]Ticket, int, error) {
	var tickets []Ticket
	var total int

	const queryTotal = `
		SELECT COUNT(*)
		FROM tickets
		WHERE status = 'OPEN'
	`

	if err := r.db.GetContext(ctx, &total, queryTotal); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, title, description, category_id, status, priority, created_by, assigned_to, assigned_at, resolved_at, closed_at, closed_by, created_at, updated_at
		FROM tickets
		WHERE status = 'OPEN'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (params.Page - 1) * params.Limit
	if err := r.db.SelectContext(ctx, &tickets, query, params.Limit, offset); err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}
