package event

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

type PasswordResetRequestedEvent struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	ResetURL  string `json:"reset_url"`
	ExpiresIn string `json:"expires_in,omitempty"`
	Forced    bool   `json:"forced,omitempty"`
}
