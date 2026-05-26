package ticket

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	CategoryID  int64      `db:"category_id"`
	DivisionID  int64      `db:"division_id"`
	Status      string     `db:"status"`
	Priority    *string    `db:"priority"`
	CreatedBy   uuid.UUID  `db:"created_by"`
	AssignedTo  *uuid.UUID `db:"assigned_to"`
	ResolvedBy  *uuid.UUID `db:"resolved_by"`
	ClosedBy    *uuid.UUID `db:"closed_by"`
	Resolution  *string    `db:"resolution"`
	AssignedAt  *time.Time `db:"assigned_at"`
	ResolvedAt  *time.Time `db:"resolved_at"`
	ClosedAt    *time.Time `db:"closed_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type TicketAttachment struct {
	ID             int64     `db:"id"`
	TicketID       int64     `db:"ticket_id"`
	FileKey        string    `db:"file_key"`
	AttachmentType string    `db:"attachment_type"`
	UploadedBy     uuid.UUID `db:"uploaded_by"`
	CreatedAt      time.Time `db:"created_at"`
}
