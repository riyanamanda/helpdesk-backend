package category

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	List(ctx context.Context, params ListCategoriesParams) ([]Category, error)
}

type repository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) List(ctx context.Context, params ListCategoriesParams) ([]Category, error) {
	var categories []Category

	const query = `
		SELECT id, name, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &categories, query, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
