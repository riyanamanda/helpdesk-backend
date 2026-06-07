package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type NotificationRepository interface {
	GetAll(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkAsRead(ctx context.Context, id int64, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	CreateBatch(ctx context.Context, notifications []Notification) error
}

type repository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, userID uuid.UUID) ([]Notification, error) {
	var notifications []Notification

	const query = `
		SELECT
			id, user_id, type, reference_type, reference_id, metadata, is_read, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 10
	`

	if err := r.db.SelectContext(ctx, &notifications, query, userID); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *repository) CreateBatch(ctx context.Context, notifications []Notification) error {
	const query = `
		INSERT INTO notifications (user_id, type, reference_type, reference_id, metadata)
		VALUES (:user_id, :type, :reference_type, :reference_id, :metadata)
	`

	for _, n := range notifications {
		if _, err := r.db.NamedExecContext(ctx, query, n); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64

	const query = `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	if err := r.db.GetContext(ctx, &count, query, userID); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	const query = `
		UPDATE notifications
		SET is_read = TRUE, read_at = NOW()
		WHERE user_id = $1 AND is_read = FALSE
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) MarkAsRead(ctx context.Context, id int64, userID uuid.UUID) error {
	const query = `
		UPDATE notifications
		SET is_read = TRUE, read_at = COALESCE(read_at, NOW())
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	return database.CheckRowsAffected(result, ErrNotificationNotFound)
}
