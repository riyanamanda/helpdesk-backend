package feedback

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type FeedbackResponse struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        FeedbackType   `json:"type"`
	Status      FeedbackStatus `json:"status"`
	CreatedBy   FeedbackUser   `json:"created_by"`
	ReviewedBy  *FeedbackUser  `json:"reviewed_by"`
	ReviewedAt  *time.Time     `json:"reviewed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type FeedbackUser struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GetFeedbackParams struct {
	pagination.Params

	Type        *string    `query:"type"`
	Status      *string    `query:"status"`
	CreatedByID *uuid.UUID `query:"-"`
}

func (p *GetFeedbackParams) Normalize() {
	p.Params.Normalize()
}

type CreateFeedbackRequest struct {
	Title       string `json:"title" validate:"required,min=5,max=100"`
	Description string `json:"description" validate:"required,min=5"`
	Type        string `json:"type" validate:"required,oneof=FEATURE_REQUEST IMPROVEMENT BUG_REPORT"`
}

type UpdateFeedbackStatusRequest struct {
	Status FeedbackStatus `json:"status" validate:"required,oneof=IN_REVIEW ACCEPTED REJECTED DELIVERED"`
}
