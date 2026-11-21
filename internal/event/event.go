package event

import (
	"time"

	"github.com/google/uuid"
)

type UserCreatedEvent struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	SetPasswordURL string `json:"set_password_url,omitempty"`
	ExpiresIn      string `json:"expires_in,omitempty"`
}

type TicketCreatedEvent struct {
	TicketID    int64  `json:"ticket_id"`
	SubmittedBy string `json:"submitted_by"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TicketAssignedEvent struct {
	TicketID     int64     `json:"ticket_id"`
	Title        string    `json:"title"`
	Priority     string    `json:"priority"`
	AssigneeID   uuid.UUID `json:"assignee_id"`
	AssignedByID uuid.UUID `json:"assigned_by_id"`
	Note         string    `json:"note,omitempty"`
}

type TicketResolvedEvent struct {
	TicketID     int64     `json:"ticket_id"`
	Title        string    `json:"title"`
	ReporterID   uuid.UUID `json:"reporter_id"`
	ResolvedByID uuid.UUID `json:"resolved_by_id"`
	Resolution   string    `json:"resolution"`
}

type PasswordChangedEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ChangedAt time.Time `json:"changed_at"`
	ByAdmin   bool      `json:"by_admin"`
}

type PasswordResetRequestedEvent struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	ResetURL  string `json:"reset_url"`
	ExpiresIn string `json:"expires_in,omitempty"`
}

type FeedbackCreatedEvent struct {
	FeedbackID    int64     `json:"feedback_id"`
	Title         string    `json:"title"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	SubmittedByID uuid.UUID `json:"submitted_by_id"`
}

type FeedbackStatusChangedEvent struct {
	FeedbackID   int64     `json:"feedback_id"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	SubmitterID  uuid.UUID `json:"submitter_id"`
	ReviewedByID uuid.UUID `json:"reviewed_by_id"`
}
