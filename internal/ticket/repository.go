package ticket

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type TicketRepository interface {
	GetAll(ctx context.Context, params GetTicketParams) ([]TicketProjection, int64, error)
	GetByID(ctx context.Context, id int64) (*TicketProjection, error)
	GetAttachmentsByTicketID(ctx context.Context, ticketID int64) (*[]TicketAttachmentProjection, error)
	Assign(ctx context.Context, ticketID int64, userID uuid.UUID) error
	UpdatePriority(ctx context.Context, ticketID int64, priority TicketPriority) error
	Update(ctx context.Context, ticketID int64, ticket Ticket) error
	CloseTicket(ctx context.Context, ticketID int64, userID uuid.UUID) error
	Begin(ctx context.Context) (TicketTx, error)
}

type TicketTx interface {
	Create(ctx context.Context, ticket Ticket) (int64, error)
	CreateAttachment(ctx context.Context, attachment TicketAttachment) error
	UpdateResolution(ctx context.Context, ticketID int64, userID uuid.UUID, resolution string) error
	DeleteAttachmentsByTicketID(ctx context.Context, ticketID int64) error
	Delete(ctx context.Context, ticketID int64) error
	Commit() error
	Rollback() error
}

type repository struct {
	db *sqlx.DB
}

type txRepository struct {
	tx *sqlx.Tx
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Begin(ctx context.Context) (TicketTx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &txRepository{tx: tx}, nil
}

func (t *txRepository) Commit() error   { return t.tx.Commit() }
func (t *txRepository) Rollback() error { return t.tx.Rollback() }

func (t *txRepository) Create(ctx context.Context, ticket Ticket) (int64, error) {
	const query = `
		INSERT INTO tickets (title, description, category_id, division_id, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := t.tx.QueryRowxContext(ctx, query, ticket.Title, ticket.Description, ticket.CategoryID, ticket.DivisionID, ticket.CreatedBy).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *txRepository) CreateAttachment(ctx context.Context, attachment TicketAttachment) error {
	const query = `
		INSERT INTO ticket_attachments (ticket_id, file_key, attachment_type, uploaded_by)
		VALUES ($1, $2, $3, $4)
	`

	_, err := t.tx.ExecContext(ctx, query, attachment.TicketID, attachment.FileKey, attachment.AttachmentType, attachment.UploadedBy)

	return err
}

func (t *txRepository) UpdateResolution(ctx context.Context, ticketID int64, userID uuid.UUID, resolution string) error {
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

	result, err := t.tx.ExecContext(ctx, query, ticketID, userID, resolution)
	if err != nil {
		return err
	}

	return database.CheckRowsAffected(result, ErrTicketNotFound)
}

func (r *repository) GetAll(ctx context.Context, params GetTicketParams) ([]TicketProjection, int64, error) {
	var (
		tickets []TicketProjection
		total   int64
	)

	where, args := buildTicketWhere(params)

	queryTotal := fmt.Sprintf(`SELECT COUNT(*) FROM tickets t %s`, where)
	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	col, dir := buildTicketSort(params)

	query := fmt.Sprintf(ticketSelectBase+`
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, col, dir, len(args)-1, len(args))
	if err := r.db.SelectContext(ctx, &tickets, query, args...); err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*TicketProjection, error) {
	var ticket TicketProjection

	const query = ticketSelectBase + `WHERE t.id = $1`

	if err := r.db.GetContext(ctx, &ticket, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

func (r *repository) GetAttachmentsByTicketID(ctx context.Context, ticketID int64) (*[]TicketAttachmentProjection, error) {
	var attachment []TicketAttachmentProjection

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
	`

	if err := r.db.SelectContext(ctx, &attachment, query, ticketID); err != nil {
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

	return database.CheckRowsAffected(result, ErrTicketNotFound)
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

	return database.CheckRowsAffected(result, ErrTicketNotFound)
}

func (t *txRepository) DeleteAttachmentsByTicketID(ctx context.Context, ticketID int64) error {
	const query = `DELETE FROM ticket_attachments WHERE ticket_id = $1`
	_, err := t.tx.ExecContext(ctx, query, ticketID)
	return err
}

func (t *txRepository) Delete(ctx context.Context, ticketID int64) error {
	const query = `DELETE FROM tickets WHERE id = $1`
	result, err := t.tx.ExecContext(ctx, query, ticketID)
	if err != nil {
		return err
	}
	return database.CheckRowsAffected(result, ErrTicketNotFound)
}

func (r *repository) Update(ctx context.Context, ticketID int64, ticket Ticket) error {
	const query = `
		UPDATE tickets
		SET title       = $2,
			description = $3,
			category_id = $4,
			division_id = $5,
			updated_at  = NOW()
		WHERE id = $1
		AND created_by = $6
		AND status = 'OPEN'
	`

	result, err := r.db.ExecContext(ctx, query,
		ticketID, ticket.Title, ticket.Description, ticket.CategoryID, ticket.DivisionID, ticket.CreatedBy,
	)
	if err != nil {
		return err
	}

	return database.CheckRowsAffected(result, ErrTicketNotFound)
}

func (r *repository) CloseTicket(ctx context.Context, ticketID int64, userID uuid.UUID) error {
	const query = `
		UPDATE tickets
		SET status = 'CLOSED',
			closed_by = $2,
			closed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		AND status != 'CLOSED'
	`

	result, err := r.db.ExecContext(ctx, query, ticketID, userID)
	if err != nil {
		return err
	}

	return database.CheckRowsAffected(result, ErrTicketNotFound)
}
