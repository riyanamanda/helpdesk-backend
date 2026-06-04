package feedback

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          int64          `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Type        FeedbackType   `db:"type"`
	Status      FeedbackStatus `db:"status"`
	CreatedBy   uuid.UUID      `db:"created_by"`
	ReviewedBy  *uuid.UUID     `db:"reviewed_by"`
	ReviewedAt  *time.Time     `db:"reviewed_at"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}
