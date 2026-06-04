package feedback

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type FeedbackRepository interface {
	GetAll(ctx context.Context, params GetFeedbackParams) ([]FeedbackProjection, int64, error)
	GetByID(ctx context.Context, id int64) (*FeedbackProjection, error)
	Create(ctx context.Context, feedback Feedback) error
	UpdateStatus(ctx context.Context, id int64, reviewerID uuid.UUID, status FeedbackStatus) error
}

type repository struct {
	db *sqlx.DB
}

func NewFeedbackRepository(db *sqlx.DB) FeedbackRepository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, params GetFeedbackParams) ([]FeedbackProjection, int64, error) {
	var (
		feedbacks []FeedbackProjection
		total     int64
		args      []any
	)

	where := "WHERE 1=1"

	if params.Type != nil {
		args = append(args, *params.Type)
		where += fmt.Sprintf(" AND f.type = $%d", len(args))
	}

	if params.Status != nil {
		args = append(args, *params.Status)
		where += fmt.Sprintf(" AND f.status = $%d", len(args))
	}

	if params.CreatedByID != nil {
		args = append(args, *params.CreatedByID)
		where += fmt.Sprintf(" AND f.created_by = $%d", len(args))
	}

	queryTotal := fmt.Sprintf(`SELECT COUNT(*) FROM feedbacks f %s`, where)
	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	sortCols := map[string]string{
		"created_at": "f.created_at",
		"updated_at": "f.updated_at",
		"status":     "f.status",
		"type":       "f.type",
	}

	col, ok := sortCols[params.SortBy]
	if !ok {
		col = "f.created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT
			f.id,
			f.title,
			f.description,
			f.type,
			f.status,
			uc.id   AS created_by_id,
			uc.name AS created_by_name,
			ur.id   AS reviewed_by_id,
			ur.name AS reviewed_by_name,
			f.reviewed_at,
			f.created_at,
			f.updated_at
		FROM feedbacks f
		JOIN users uc ON uc.id = f.created_by
		LEFT JOIN users ur ON ur.id = f.reviewed_by
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, col, dir, len(args)-1, len(args))

	if err := r.db.SelectContext(ctx, &feedbacks, query, args...); err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*FeedbackProjection, error) {
	var feedback FeedbackProjection

	const query = `
		SELECT
			f.id,
			f.title,
			f.description,
			f.type,
			f.status,
			uc.id   AS created_by_id,
			uc.name AS created_by_name,
			ur.id   AS reviewed_by_id,
			ur.name AS reviewed_by_name,
			f.reviewed_at,
			f.created_at,
			f.updated_at
		FROM feedbacks f
		JOIN users uc ON uc.id = f.created_by
		LEFT JOIN users ur ON ur.id = f.reviewed_by
		WHERE f.id = $1
	`

	if err := r.db.GetContext(ctx, &feedback, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFeedbackNotFound
		}
		return nil, err
	}

	return &feedback, nil
}

func (r *repository) Create(ctx context.Context, feedback Feedback) error {
	const query = `
		INSERT INTO feedbacks (title, description, type, created_by)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, feedback.Title, feedback.Description, feedback.Type, feedback.CreatedBy)
	return err
}

func (r *repository) UpdateStatus(ctx context.Context, id int64, reviewerID uuid.UUID, status FeedbackStatus) error {
	const query = `
		UPDATE feedbacks
		SET status      = $3,
		    reviewed_by = $2,
		    reviewed_at = NOW(),
		    updated_at  = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, reviewerID, status)
	if err != nil {
		return err
	}

	return database.CheckRowsAffected(result, ErrFeedbackNotFound)
}
