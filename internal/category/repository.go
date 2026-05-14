package category

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

//go:generate mockery --name CategoryRepository
type CategoryRepository interface {
	List(ctx context.Context, params GetCategoryParams) ([]Category, int, error)
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

func (r *repository) List(ctx context.Context, params GetCategoryParams) ([]Category, int, error) {
	var categories []Category
	var total int

	const queryTotal = `
		SELECT COUNT(*)
		FROM categories
		WHERE is_active = TRUE
	`

	if err := r.db.GetContext(ctx, &total, queryTotal); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, name, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (params.Page - 1) * params.Limit
	if err := r.db.SelectContext(ctx, &categories, query, params.Limit, offset); err != nil {
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
		WHERE id = $1 AND is_active = TRUE
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
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND is_active = TRUE
	`

	result, err := r.db.ExecContext(ctx, query, category.Name, id)

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
