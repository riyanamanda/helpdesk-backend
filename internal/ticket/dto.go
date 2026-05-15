package ticket

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type TicketResponse struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CategoryID  int             `json:"category_id"`
	Status      TicketStatus    `json:"status"`
	Priority    *TicketPriority `json:"priority"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	AssignedTo  *uuid.UUID      `json:"assigned_to"`
	AssignedAt  *time.Time      `json:"assigned_at"`
	ResolvedAt  *time.Time      `json:"resolved_at"`
	ClosedAt    *time.Time      `json:"closed_at"`
	ClosedBy    *uuid.UUID      `json:"closed_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type GetTicketParams struct {
	pagination.Params
}

type TicketCreateRequest struct {
	Title       string `form:"title" validate:"required,min=5,max=100"`
	Description string `form:"description" validate:"required,min=5,max=255"`
	CategoryID  int    `form:"category_id" validate:"required,gt=0"`
}
