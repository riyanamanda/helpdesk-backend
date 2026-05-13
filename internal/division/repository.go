package division

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	dberror "github.com/riyanamanda/helpdesk-backend/internal/infra/database"
)

//go:generate mockery --name DivisionRepository
type DivisionRepository interface {
	List(ctx context.Context, params GetDivisionParams) ([]Division, int, error)
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

func (r *repository) List(ctx context.Context, params GetDivisionParams) ([]Division, int, error) {
	var divisions []Division
	var total int

	const queryTotal = `
		SELECT COUNT(*)
		FROM divisions
		WHERE is_active = TRUE
	`
	if err := r.db.GetContext(ctx, &total, queryTotal); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, name, is_active, created_at, updated_at
		FROM divisions
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	offset := (params.Page - 1) * params.Limit
	if err := r.db.SelectContext(ctx, &divisions, query, params.Limit, offset); err != nil {
		return nil, 0, err
	}

	return divisions, total, nil
}

func (r *repository) Create(ctx context.Context, division *Division) error {
	const query = `
		INSERT INTO divisions (name)
		VALUES ($1)
		RETURNING id, name, is_active, created_at, updated_at
	`
	err := r.db.QueryRowxContext(ctx, query, division.Name).
		Scan(
			&division.ID,
			&division.Name,
			&division.IsActive,
			&division.CreatedAt,
			&division.UpdatedAt,
		)
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
		WHERE id = $1 AND is_active = TRUE
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
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND is_active = TRUE
		RETURNING id, name, is_active, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, division.Name, id).
		Scan(
			&division.ID,
			&division.Name,
			&division.IsActive,
			&division.CreatedAt,
			&division.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDivisionNotFound
		}
		if dberror.IsUniqueViolation(err) {
			return ErrDivisionAlreadyExists
		}
		return err
	}

	return nil
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

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrDivisionNotFound
	}

	return nil
}
