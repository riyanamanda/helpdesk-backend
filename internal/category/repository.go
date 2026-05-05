package category

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	List(ctx context.Context, params ListCategoriesParams) ([]Category, int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) List(ctx context.Context, params ListCategoriesParams) ([]Category, int, error) {
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

	if err := r.db.SelectContext(ctx, &categories, query, params.Limit, params.Offset); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
