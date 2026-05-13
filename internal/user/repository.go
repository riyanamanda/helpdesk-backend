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
	List(ctx context.Context, params GetUserParams) ([]User, int, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) List(ctx context.Context, params GetUserParams) ([]User, int, error) {
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
		SELECT id, name, email, avatar_key, phone, role, division_id, is_active, created_by, created_at, updated_at
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
		RETURNING id, name, email, role, division_id, is_active, created_by, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, user.Name, user.Email, user.Password, user.Role, user.DivisionID, user.CreatedBy).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)

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
		SELECT id, name, email, avatar_key, phone, role, division_id, is_active, created_by, created_at, updated_at
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
		SELECT id, name, email, password, avatar_key, phone, role, division_id, is_active, created_by, created_at, updated_at
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
