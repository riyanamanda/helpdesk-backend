package mailer

const QueueNewTicketEmail = "email.new_ticket"

type newTicketPayload struct {
	TicketID    int64  `json:"ticket_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SubmitterID string `json:"submitter_id"`
}
