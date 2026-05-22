package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

//go:generate mockery --name CategoryRepository
type CategoryRepository interface {
	GetAll(ctx context.Context, params GetCategoryParams) ([]Category, int, error)
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	Update(ctx context.Context, id int64, category *Category) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetCategoryParams) ([]Category, int, error) {
	var (
		categories []Category
		total      int
		args       []any
	)

	where := "WHERE 1=1"

	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")

		where += fmt.Sprintf(" AND name ILIKE $%d", len(args))
	}

	if params.IsActive != nil {
		args = append(args, *params.IsActive)

		where += fmt.Sprintf(" AND is_active = $%d", len(args))
	}

	queryTotal := fmt.Sprintf(`SELECT COUNT(*) FROM categories %s`, where)

	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			is_active,
			created_at,
			updated_at
		FROM categories
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, params.SortBy, params.SortType, len(args)-1, len(args))

	if err := r.db.SelectContext(ctx, &categories, query, args...); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *repository) Create(ctx context.Context, category *Category) error {
	const query = `
		INSERT INTO categories (name)
		VALUES ($1)
	`

	_, err := r.db.ExecContext(ctx, query, category.Name)

	if err != nil {
		if dberror.IsUniqueViolation(err) {
			return ErrCategoryAlreadyExists
		}
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Category, error) {
	var category Category

	const query = `
		SELECT id, name, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	if err := r.db.GetContext(ctx, &category, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *repository) Update(ctx context.Context, id int64, category *Category) error {
	const query = `
		UPDATE categories
		SET name = $2,
			is_active = $3,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, category.Name, category.IsActive)

	if err != nil {
		if dberror.IsUniqueViolation(err) {
			return ErrCategoryAlreadyExists
		}
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM categories
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}
