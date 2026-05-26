package ticket

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type TicketResponse struct {
	ID          int64                  `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Category    category.CategoryBrief `json:"category"`
	Division    division.DivisionBrief `json:"division"`
	Status      TicketStatus           `json:"status"`
	Priority    *TicketPriority        `json:"priority"`
	CreatedBy   user.UserBrief         `json:"created_by"`
	AssignedTo  *user.UserBrief        `json:"assigned_to"`
	ResolvedBy  *user.UserBrief        `json:"resolved_by"`
	ClosedBy    *user.UserBrief        `json:"closed_by"`
	Resolution  *string                `json:"resolution"`
	AssignedAt  *time.Time             `json:"assigned_at"`
	ResolvedAt  *time.Time             `json:"resolved_at"`
	ClosedAt    *time.Time             `json:"closed_at"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TicketDetailResponse struct {
	TicketResponse
	Attachments *[]TicketAttachmentResponse `json:"attachment"`
}

type TicketAttachmentResponse struct {
	ID             int64          `json:"id"`
	TicketID       int64          `json:"ticket_id"`
	FileURL        string         `json:"file_url"`
	AttachmentType string         `json:"attachment_type"`
	UploadedBy     user.UserBrief `json:"uploaded_by"`
	CreatedAt      time.Time      `json:"created_at"`
}

type GetTicketParams struct {
	pagination.Params
}

type TicketCreateRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=5,max=100"`
	Description string `json:"description" form:"description" validate:"required,min=5,max=255"`
	CategoryID  int64  `json:"category" form:"category" validate:"required,gt=0"`
	DivisionID  int64  `json:"division" form:"division" validate:"required,gt=0"`
}

type TicketAssignRequest struct {
	AssignedTo uuid.UUID `json:"assigned_to" validate:"required"`
}

type TicketPriorityRequest struct {
	Priority TicketPriority `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH URGENT"`
}

type TicketResolutionRequest struct {
	Resolution string `json:"resolution" form:"resolution" validate:"required,max=1000"`
}
