package division

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

type DivisionRepository interface {
	GetAll(ctx context.Context, params GetDivisionParams) ([]Division, int64, error)
	GetOptions(ctx context.Context) ([]DivisionOptionProjection, error)
	Create(ctx context.Context, division *Division) error
	GetByID(ctx context.Context, id int64) (*Division, error)
	Update(ctx context.Context, id int64, division *Division) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *sqlx.DB
}

func NewDivisionRepository(db *sqlx.DB) DivisionRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetDivisionParams) ([]Division, int64, error) {
	var (
		divisions []Division
		total     int64
		args      []any
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

	queryTotal := fmt.Sprintf(`SELECT COUNT(*) FROM divisions %s`, where)

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
		FROM divisions
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, col, dir, len(args)-1, len(args))

	if err := r.db.SelectContext(ctx, &divisions, query, args...); err != nil {
		return nil, 0, err
	}

	return divisions, total, nil
}

func (r *repository) GetOptions(ctx context.Context) ([]DivisionOptionProjection, error) {
	var divisions []DivisionOptionProjection

	const query = `
		SELECT
			id,
			name
		FROM divisions
		WHERE is_active = true
		ORDER BY name ASC
	`

	if err := r.db.SelectContext(ctx, &divisions, query); err != nil {
		return nil, err
	}

	return divisions, nil
}

func (r *repository) Create(ctx context.Context, division *Division) error {
	const query = `
		INSERT INTO divisions (name)
		VALUES ($1)
	`

	_, err := r.db.ExecContext(ctx, query, division.Name)
	if err != nil {
		if dberror.IsUniqueViolation(err) {
			return ErrDivisionAlreadyExists
		}
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Division, error) {
	var division Division

	const query = `
		SELECT id, name, is_active, created_at, updated_at
		FROM divisions
		WHERE id = $1
	`

	if err := r.db.GetContext(ctx, &division, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDivisionNotFound
		}
		return nil, err
	}

	return &division, nil
}

func (r *repository) Update(ctx context.Context, id int64, division *Division) error {
	const query = `
		UPDATE divisions
		SET name = $2,
			is_active = $3,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, division.Name, division.IsActive)
	if err != nil {
		if dberror.IsUniqueViolation(err) {
			return ErrDivisionAlreadyExists
		}
		return err
	}

	return dberror.CheckRowsAffected(result, ErrDivisionNotFound)
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM divisions
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return dberror.CheckRowsAffected(result, ErrDivisionNotFound)
}
