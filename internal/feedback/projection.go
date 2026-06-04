package feedback

import (
	"time"

	"github.com/google/uuid"
)

type FeedbackProjection struct {
	ID             int64          `db:"id"`
	Title          string         `db:"title"`
	Description    string         `db:"description"`
	Type           FeedbackType   `db:"type"`
	Status         FeedbackStatus `db:"status"`
	CreatedByID    uuid.UUID      `db:"created_by_id"`
	CreatedByName  string         `db:"created_by_name"`
	ReviewedByID   *uuid.UUID     `db:"reviewed_by_id"`
	ReviewedByName *string        `db:"reviewed_by_name"`
	ReviewedAt     *time.Time     `db:"reviewed_at"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}
