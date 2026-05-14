package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	GetAll(ctx context.Context, params GetUserParams) ([]User, int, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarKey string) error
}

type repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetUserParams) ([]User, int, error) {
	var users []User
	var total int

	const queryTotal = `
		SELECT COUNT(*)
		FROM users
		WHERE is_active = TRUE
	`

	if err := r.db.GetContext(ctx, &total, queryTotal); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, name, email, google_id, avatar_key, phone, role, division_id, is_active, created_by, created_at, updated_at
		FROM users
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (params.Page - 1) * params.Limit
	if err := r.db.SelectContext(ctx, &users, query, params.Limit, offset); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) Create(ctx context.Context, user *User) error {
	const query = `
		INSERT INTO users (name, email, password, role, division_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role, user.DivisionID, user.CreatedBy)

	if err != nil {
		if dberror.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User

	const query = `
		SELECT id, name, email, google_id, avatar_key, phone, role, division_id, is_active, created_by, created_at, updated_at
		FROM users
		WHERE id = $1
		AND is_active = TRUE
	`

	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	const query = `
		SELECT id, name, email, password, google_id, avatar_key, phone, role, division_id, is_active, created_by, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`

	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
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

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
