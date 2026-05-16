package ticket

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type TicketResponse struct {
	ID          int64                  `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Category    category.CategoryBrief `json:"category"`
	Status      TicketStatus           `json:"status"`
	Priority    *TicketPriority        `json:"priority"`
	CreatedBy   user.UserBrief         `json:"created_by"`
	AssignedTo  *user.UserBrief        `json:"assigned_to"`
	AssignedAt  *time.Time             `json:"assigned_at"`
	ResolvedAt  *time.Time             `json:"resolved_at"`
	ClosedAt    *time.Time             `json:"closed_at"`
	ClosedBy    *user.UserBrief        `json:"closed_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TicketDetailResponse struct {
	TicketResponse
	Attachment *TicketAttachmentResponse `json:"attachment"`
}

type TicketAttachmentResponse struct {
	ID             int64          `json:"id"`
	TicketID       int64          `json:"ticket_id"`
	FileURL        string         `json:"file_url"`
	AttachmentType string         `json:"attachment_type"`
	UploadedBy     user.UserBrief `json:"uploaded_by"`
	CreatedAt      time.Time      `json:"created_at"`
}

type TicketResolutionResponse struct {
	ID         int64          `json:"id"`
	TicketID   int64          `json:"ticket_id"`
	ResolvedBy user.UserBrief `json:"resolved_by"`
	Resolution string         `json:"resolution"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type GetTicketParams struct {
	pagination.Params
}

type TicketCreateRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=5,max=100"`
	Description string `json:"description" form:"description" validate:"required,min=5,max=255"`
	CategoryID  int    `json:"category_id" form:"category_id" validate:"required,gt=0"`
}

type TicketAssignRequest struct {
	AssignedTo uuid.UUID `json:"assigned_to" validate:"required"`
}

type TicketResolutionCreateRequest struct {
	Resolution string `json:"resolution" form:"resolution" validate:"required,max=1000"`
}
