package user_device

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type UserDeviceRepository interface {
	Upsert(ctx context.Context, userID uuid.UUID, fcmToken string) error
	Delete(ctx context.Context, userID uuid.UUID, fcmToken string) error
	GetTokensByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]string, error)
}

type repository struct {
	db *sqlx.DB
}

func NewUserDeviceRepository(db *sqlx.DB) UserDeviceRepository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, userID uuid.UUID, fcmToken string) error {
	const query = `
		INSERT INTO user_devices (user_id, fcm_token, last_seen_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (fcm_token) DO UPDATE
		SET user_id = EXCLUDED.user_id, last_seen_at = NOW(), updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, userID, fcmToken)
	return err
}

func (r *repository) Delete(ctx context.Context, userID uuid.UUID, fcmToken string) error {
	const query = `
		DELETE FROM user_devices
		WHERE user_id = $1 AND fcm_token = $2
	`
	result, err := r.db.ExecContext(ctx, query, userID, fcmToken)
	if err != nil {
		return err
	}
	return database.CheckRowsAffected(result, ErrDeviceNotFound)
}

func (r *repository) GetTokensByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var tokens []string
	const query = `SELECT fcm_token FROM user_devices WHERE user_id = $1`
	if err := r.db.SelectContext(ctx, &tokens, query, userID); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *repository) GetTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]string, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`SELECT fcm_token FROM user_devices WHERE user_id IN (?)`, userIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var tokens []string
	if err := r.db.SelectContext(ctx, &tokens, query, args...); err != nil {
		return nil, err
	}
	return tokens, nil
}
