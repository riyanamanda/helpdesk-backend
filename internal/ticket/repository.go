package ticket

import (
	"context"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name TicketRepository
type TicketRepository interface {
	GetAll(ctx context.Context, params GetTicketParams) ([]TicketProjection, int64, error)
	Create(ctx context.Context, ticket Ticket) (int64, error)
	CreateAttachment(ctx context.Context, attachment TicketAttachment) error
}

type repository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetTicketParams) ([]TicketProjection, int64, error) {
	var tickets []TicketProjection
	var total int64

	const queryTotal = `
		SELECT COUNT(*)
		FROM tickets
		WHERE status = 'OPEN'
	`

	if err := r.db.GetContext(ctx, &total, queryTotal); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT
			t.id,
			t.title,
			t.description,
			c.id AS category_id,
			c.name AS category_name,
			t.status,
			t.priority,
			u.id AS created_by_id,
			u.name AS created_by_name,
			uat.id AS assigned_to_id,
			uat.name AS assigned_to_name,
			t.assigned_at,
			t.resolved_at,
			t.closed_at,
			ucb.id AS closed_by_id,
			ucb.name AS closed_by_name,
			t.created_at,
			t.updated_at
		FROM tickets t
		JOIN categories c
			ON c.id = t.category_id
		JOIN users u
			ON u.id = t.created_by
		LEFT JOIN users uat
			ON uat.id = t.assigned_to
		LEFT JOIN users ucb
			ON ucb.id = t.closed_by
		WHERE t.status = 'OPEN'
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (params.Page - 1) * params.Limit
	if err := r.db.SelectContext(ctx, &tickets, query, params.Limit, offset); err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

func (r *repository) Create(ctx context.Context, ticket Ticket) (int64, error) {
	const query = `
		INSERT INTO tickets (title, description, category_id, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowxContext(ctx, query, ticket.Title, ticket.Description, ticket.CategoryID, ticket.CreatedBy).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) CreateAttachment(ctx context.Context, attachment TicketAttachment) error {
	const query = `
		INSERT INTO ticket_attachments (ticket_id, file_key, attachment_type, uploaded_by)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, attachment.TicketID, attachment.FileKey, attachment.AttachmentType, attachment.UploadedBy)
	if err != nil {
		return err
	}

	return nil
}
