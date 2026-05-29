package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

type UserRepository interface {
	GetAll(ctx context.Context, params GetUserParams) ([]UserProjection, int64, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*UserProjection, error)
	GetByEmail(ctx context.Context, email string) (*UserProjection, error)
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
		users []UserProjection
		total int64
		args  []any
	)

	where := "WHERE 1=1"

	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		where += fmt.Sprintf(" AND u.name ILIKE $%d", len(args))
	}

	if params.IsActive != nil {
		args = append(args, *params.IsActive)
		where += fmt.Sprintf(" AND u.is_active = $%d", len(args))
	}

	if params.Role != "" {
		args = append(args, params.Role)
		where += fmt.Sprintf(" AND u.role = $%d", len(args))
	}

	if params.DivisionID != nil {
		args = append(args, *params.DivisionID)
		where += fmt.Sprintf(" AND u.division_id = $%d", len(args))
	}

	queryTotal := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, where)
	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)
	sortCols := map[string]string{
		"name": "u.name", "role": "u.role", "division_id": "u.division_id",
		"is_active": "u.is_active", "created_at": "u.created_at",
	}

	col, ok := sortCols[params.SortBy]
	if !ok {
		col = "u.created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.name,
			u.email,
			u.google_id,
			u.avatar_key,
			u.phone,
			u.role,
			u.gender,
			d.id as division_id,
			d.name as division_name,
			u.is_active,
			cb.id as created_by_id,
			cb.name as created_by_name,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN divisions d
			ON d.id = u.division_id
		LEFT JOIN users cb
			ON cb.id = u.created_by
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
		if dberror.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}

		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*UserProjection, error) {
	var user UserProjection

	const query = `
		SELECT
			u.id,
			u.name,
			u.email,
			u.google_id,
			u.avatar_key,
			u.phone,
			u.role,
			u.gender,
			d.id as division_id,
			d.name as division_name,
			u.is_active,
			cb.id as created_by_id,
			cb.name as created_by_name,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN divisions d
			ON d.id = u.division_id
		LEFT JOIN users cb
			ON cb.id = u.created_by
		WHERE u.id = $1
	`

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

	const query = `
		SELECT
			u.id,
			u.name,
			u.email,
			u.password,
			u.google_id,
			u.avatar_key,
			u.phone,
			u.role,
			u.gender,
			d.id as division_id,
			d.name as division_name,
			u.is_active,
			cb.id as created_by_id,
			cb.name as created_by_name,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN divisions d
			ON d.id = u.division_id
		LEFT JOIN users cb
			ON cb.id = u.created_by
		WHERE LOWER(u.email) = LOWER($1)
	`

	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
