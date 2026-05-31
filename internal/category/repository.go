package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

type CategoryRepository interface {
	GetAll(ctx context.Context, params GetCategoryParams) ([]Category, int64, error)
	GetOptions(ctx context.Context) ([]CategoryOptionProjection, error)
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

func (r *repository) GetAll(ctx context.Context, params GetCategoryParams) ([]Category, int64, error) {
	var (
		categories []Category
		total      int64
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
	sortCols := map[string]string{
		"name": "name", "is_active": "is_active", "created_at": "created_at",
	}

	col, ok := sortCols[params.SortBy]
	if !ok {
		col = "created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

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
	`, where, col, dir, len(args)-1, len(args))

	if err := r.db.SelectContext(ctx, &categories, query, args...); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *repository) GetOptions(ctx context.Context) ([]CategoryOptionProjection, error) {
	var categories []CategoryOptionProjection

	const query = `
		SELECT
			id,
			name
		FROM categories
		WHERE is_active = true
		ORDER BY name ASC
	`

	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, err
	}

	return categories, nil
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

	return dberror.CheckRowsAffected(result, ErrCategoryNotFound)
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

	return dberror.CheckRowsAffected(result, ErrCategoryNotFound)
}
