package category

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CategoryRepository interface {
	List(ctx context.Context, params GetCategoryParams) ([]Category, int, error)
	Create(ctx context.Context, category *Category) error
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

	_, limit, offset := params.Normalize()
	if err := r.db.SelectContext(ctx, &categories, query, limit, offset); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *repository) Create(ctx context.Context, category *Category) error {
	const query = `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id, name, is_active, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, category.Name).
		Scan(
			&category.ID,
			&category.Name,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrCategoryAlreadyExists
		}
		return err
	}

	return nil
}
