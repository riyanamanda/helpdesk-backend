package ticket

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/database"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

//go:generate mockery --name TicketRepository
type TicketRepository interface {
	GetAll(ctx context.Context, params GetTicketParams) ([]TicketProjection, int64, error)
	Create(ctx context.Context, tx *sqlx.Tx, ticket Ticket) (int64, error)
	GetByID(ctx context.Context, id int64) (*TicketProjection, error)

	CreateAttachment(ctx context.Context, tx *sqlx.Tx, attachment TicketAttachment) error
	GetAttachmentByTicketID(ctx context.Context, ticketID int64, attachmentType AttachmentType) (*TicketAttachmentProjection, error)

	Assign(ctx context.Context, ticketID int64, userID uuid.UUID) error
	UpdatePriority(ctx context.Context, ticketID int64, priority TicketPriority) error
	UpdateResolution(ctx context.Context, tx *sqlx.Tx, ticketID int64, userID uuid.UUID, resolution string) error
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
		WHERE status != 'CLOSED'
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
			urb.id AS resolved_by_id,
			urb.name AS resolved_by_name,
			ucb.id AS closed_by_id,
			ucb.name AS closed_by_name,
			t.resolution,
			t.assigned_at,
			t.resolved_at,
			t.closed_at,
			t.created_at,
			t.updated_at
		FROM tickets t
		JOIN categories c
			ON c.id = t.category_id
		JOIN users u
			ON u.id = t.created_by
		LEFT JOIN users uat
			ON uat.id = t.assigned_to
		LEFT JOIN users urb
			ON urb.id = t.resolved_by
		LEFT JOIN users ucb
			ON ucb.id = t.closed_by
		WHERE t.status != 'CLOSED'
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (params.Page - 1) * params.Limit
	if err := r.db.SelectContext(ctx, &tickets, query, params.Limit, offset); err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, ticket Ticket) (int64, error) {
	const query = `
		INSERT INTO tickets (title, description, category_id, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowxContext(ctx, query, ticket.Title, ticket.Description, ticket.CategoryID, ticket.CreatedBy).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*TicketProjection, error) {
	var ticket TicketProjection

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
			urb.id AS resolved_by_id,
			urb.name AS resolved_by_name,
			ucb.id AS closed_by_id,
			ucb.name AS closed_by_name,
			t.resolution,
			t.assigned_at,
			t.resolved_at,
			t.closed_at,
			t.created_at,
			t.updated_at
		FROM tickets t
		JOIN categories c
			ON c.id = t.category_id
		JOIN users u
			ON u.id = t.created_by
		LEFT JOIN users uat
			ON uat.id = t.assigned_to
		LEFT JOIN users urb
			ON urb.id = t.resolved_by
		LEFT JOIN users ucb
			ON ucb.id = t.closed_by
		WHERE t.id = $1
	`

	if err := r.db.GetContext(ctx, &ticket, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

func (r *repository) CreateAttachment(ctx context.Context, tx *sqlx.Tx, attachment TicketAttachment) error {
	const query = `
		INSERT INTO ticket_attachments (ticket_id, file_key, attachment_type, uploaded_by)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.ExecContext(ctx, query, attachment.TicketID, attachment.FileKey, attachment.AttachmentType, attachment.UploadedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAttachmentByTicketID(ctx context.Context, ticketID int64, attachmentType AttachmentType) (*TicketAttachmentProjection, error) {
	var attachment TicketAttachmentProjection

	const query = `
		SELECT
			a.id,
			a.ticket_id,
			a.file_key,
			a.attachment_type,
			au.id AS uploaded_by_id,
			au.name AS uploaded_by_name,
			a.created_at
		FROM ticket_attachments a
		JOIN users au
			ON au.id = a.uploaded_by
		WHERE a.ticket_id = $1
		AND attachment_type = $2
		LIMIT 1
	`

	if err := r.db.GetContext(ctx, &attachment, query, ticketID, attachmentType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &attachment, nil
}

func (r *repository) Assign(ctx context.Context, ticketID int64, userID uuid.UUID) error {
	const query = `
		UPDATE tickets
		SET assigned_to = $2,
			assigned_at = NOW(),
			status = 'IN_PROGRESS',
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, ticketID, userID)
	if err != nil {
		if database.IsForeignKeyViolation(err) {
			return user.ErrUserNotFound
		}
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrTicketNotFound
	}

	return nil
}

func (r *repository) UpdatePriority(ctx context.Context, ticketID int64, priority TicketPriority) error {
	const query = `
		UPDATE tickets
		SET priority = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, ticketID, priority)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrTicketNotFound
	}

	return nil
}

func (r *repository) UpdateResolution(ctx context.Context, tx *sqlx.Tx, ticketID int64, userID uuid.UUID, resolution string) error {
	const query = `
		UPDATE tickets
		SET resolution = $3,
			resolved_by = $2,
			resolved_at = NOW(),
			status = 'RESOLVED',
			updated_at = NOW()
		WHERE id = $1
		AND status not in ('RESOLVED','CLOSED')
	`

	result, err := tx.ExecContext(ctx, query, ticketID, userID, resolution)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrTicketNotFound
	}

	return nil
}
