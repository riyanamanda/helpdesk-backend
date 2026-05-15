package ticket

import (
	"context"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name TicketRepository
type TicketRepository interface {
	GetAll(ctx context.Context, params GetTicketParams) ([]Ticket, int64, error)
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

func (r *repository) GetAll(ctx context.Context, params GetTicketParams) ([]Ticket, int64, error) {
	var tickets []Ticket
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
