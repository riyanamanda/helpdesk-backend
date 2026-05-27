package profile

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

//go:generate mockery --name ProfileRepository
type ProfileRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.UserProjection, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, name string, phone *string) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarKey string) error
	SetGoogleID(ctx context.Context, id uuid.UUID, googleID string) error
}

type repository struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*user.UserProjection, error) {
	var p user.UserProjection

	const query = `
		SELECT
			u.id,
			u.name,
			u.email,
			u.google_id,
			u.avatar_key,
			u.phone,
			u.role,
			d.id   AS division_id,
			d.name AS division_name,
			u.is_active,
			cb.id   AS created_by_id,
			cb.name AS created_by_name,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN divisions d  ON d.id  = u.division_id
		LEFT JOIN users    cb ON cb.id = u.created_by
		WHERE u.id = $1
	`

	if err := r.db.GetContext(ctx, &p, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *repository) UpdateProfile(ctx context.Context, id uuid.UUID, name string, phone *string) error {
	const query = `
		UPDATE users
		SET name       = $2,
		    phone      = $3,
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, name, phone)
	if err != nil {
		return err
	}

	return dberror.CheckRowsAffected(result, ErrProfileNotFound)
}

func (r *repository) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarKey string) error {
	const query = `
		UPDATE users
		SET avatar_key = $2,
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, avatarKey)
	if err != nil {
		return err
	}

	return dberror.CheckRowsAffected(result, ErrProfileNotFound)
}

func (r *repository) SetGoogleID(ctx context.Context, id uuid.UUID, googleID string) error {
	const query = `
		UPDATE users
		SET google_id  = $2,
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, googleID)
	if err != nil {
		if dberror.IsUniqueViolation(err) {
			return ErrGoogleIDAlreadyLinked
		}
		return err
	}

	return dberror.CheckRowsAffected(result, ErrProfileNotFound)
}
