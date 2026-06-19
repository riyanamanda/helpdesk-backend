package mailer

const (
	QueueNewTicketEmail   = "email.new_ticket"
	QueueWelcomeUserEmail = "email.welcome_user"
)

type newTicketPayload struct {
	TicketID    int64  `json:"ticket_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SubmitterID string `json:"submitter_id"`
}

type welcomeUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
