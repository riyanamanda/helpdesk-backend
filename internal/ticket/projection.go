package ticket

import (
	"time"

	"github.com/google/uuid"
)

type TicketProjection struct {
	ID int64 `db:"id"`

	Title string `db:"title"`

	Description string `db:"description"`

	CategoryID int64 `db:"category_id"`

	CategoryName string `db:"category_name"`

	DivisionID int64 `db:"division_id"`

	DivisionName string `db:"division_name"`

	Status string `db:"status"`

	Priority *string `db:"priority"`

	CreatedByID uuid.UUID `db:"created_by_id"`

	CreatedByName string `db:"created_by_name"`

	AssignedToID *uuid.UUID `db:"assigned_to_id"`

	AssignedToName *string `db:"assigned_to_name"`

	ResolvedByID *uuid.UUID `db:"resolved_by_id"`

	ResolvedByName *string `db:"resolved_by_name"`

	ClosedByID *uuid.UUID `db:"closed_by_id"`

	ClosedByName *string `db:"closed_by_name"`

	Resolution *string `db:"resolution"`

	AssignedAt *time.Time `db:"assigned_at"`

	ResolvedAt *time.Time `db:"resolved_at"`

	ClosedAt *time.Time `db:"closed_at"`

	CreatedAt time.Time `db:"created_at"`

	UpdatedAt time.Time `db:"updated_at"`
}

type TicketAttachmentProjection struct {
	ID int64 `db:"id"`

	TicketID int64 `db:"ticket_id"`

	FileKey string `db:"file_key"`

	AttachmentType string `db:"attachment_type"`

	UploadedByID uuid.UUID `db:"uploaded_by_id"`

	UploadedByName string `db:"uploaded_by_name"`

	CreatedAt time.Time `db:"created_at"`
}
