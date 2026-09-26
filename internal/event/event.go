package event

type UserCreatedEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TicketCreatedEvent struct {
	TicketID    int64  `json:"ticket_id"`
	SubmittedBy string `json:"submitted_by"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
