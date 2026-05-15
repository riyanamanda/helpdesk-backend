package ticket

import (
	"time"

	"github.com/google/uuid"
)

type TicketProjection struct {
	ID             int64      `db:"id"`
	Title          string     `db:"title"`
	Description    string     `db:"description"`
	CategoryID     int        `db:"category_id"`
	CategoryName   string     `db:"category_name"`
	Status         string     `db:"status"`
	Priority       *string    `db:"priority"`
	CreatedByID    uuid.UUID  `db:"created_by_id"`
	CreatedByName  string     `db:"created_by_name"`
	AssignedToID   *uuid.UUID `db:"assigned_to_id"`
	AssignedToName *string    `db:"assigned_to_name"`
	AssignedAt     *time.Time `db:"assigned_at"`
	ResolvedAt     *time.Time `db:"resolved_at"`
	ClosedAt       *time.Time `db:"closed_at"`
	ClosedByID     *uuid.UUID `db:"closed_by_id"`
	ClosedByName   *string    `db:"closed_by_name"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}
