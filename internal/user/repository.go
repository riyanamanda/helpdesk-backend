package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type UserRepository interface {
	GetAll(ctx context.Context, params GetUserParams) ([]UserProjection, int64, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*UserProjection, error)
	GetByEmail(ctx context.Context, email string) (*UserProjection, error)
	UpdateByID(ctx context.Context, id uuid.UUID, user User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	AssignableUser(ctx context.Context) ([]AssignableUserProjection, error)
	GetEmailsByRole(ctx context.Context, role UserRole) ([]string, error)
}

type repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetUserParams) ([]UserProjection, int64, error) {
	var (
		total int64
		users []UserProjection
	)

	where, args := buildUserWhere(params)

	queryTotal := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, where)
	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	col, dir := buildUserSort(params)

	query := fmt.Sprintf(userSelectBase+`
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, col, dir, len(args)-1, len(args))
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) Create(ctx context.Context, user *User) error {
	const query = `
		INSERT INTO users (name, email, password, role, gender, division_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role, user.Gender, user.DivisionID, user.CreatedBy)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}

		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*UserProjection, error) {
	var user UserProjection

	const query = userSelectBase + `WHERE u.id = $1`

	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*UserProjection, error) {
	var user UserProjection

	const query = userSelectWithPassword + `WHERE LOWER(u.email) = LOWER($1)`

	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *repository) UpdateByID(ctx context.Context, id uuid.UUID, user User) error {
	const query = `
		UPDATE users
		SET name 		= $2,
			email 		= $3,
			role 		= $4,
			division_id	= $5,
			gender 		= $6,
			is_active	= $7,
			updated_at	= NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, user.Name, user.Email, user.Role, user.DivisionID, user.Gender, user.IsActive)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return database.CheckRowsAffected(result, ErrUserNotFound)
}

func (r *repository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	const query = `
		UPDATE users
		SET password	= $2,
			updated_at	= NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, password)
	if err != nil {
		return err
	}

	return database.CheckRowsAffected(result, ErrUserNotFound)
}

func (r *repository) GetEmailsByRole(ctx context.Context, role UserRole) ([]string, error) {
	var emails []string

	const query = `
		SELECT email
		FROM users
		WHERE role = $1
		AND is_active = true
	`

	if err := r.db.SelectContext(ctx, &emails, query, role); err != nil {
		return nil, err
	}
	return emails, nil
}

func (r *repository) AssignableUser(ctx context.Context) ([]AssignableUserProjection, error) {
	var users []AssignableUserProjection

	const query = `
		SELECT
			id,
			name
		FROM users
		WHERE is_active = true
		AND division_id = 1
		ORDER BY name ASC
	`

	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}

	return users, nil
}
