package rabbitmq

const (
	// exchange
	ExchangeEvent      = "helpdesk.events"
	ExchangeWelcomeDLX = "helpdesk.welcome.dlx"
	ExchangeTicketDLX  = "helpdesk.ticket.dlx"

	ExchangeFeedbackDLX = "helpdesk.feedback.dlx"

	// queue
	QueueWelcome    = "helpdesk.welcome"
	QueueWelcomeDLQ = "helpdesk.welcome.dlq"
	QueueTicket     = "helpdesk.ticket"
	QueueTicketDLQ  = "helpdesk.ticket.dlq"

	QueueTicketActivity = "helpdesk.ticket-activity"

	QueueFeedback    = "helpdesk.feedback"
	QueueFeedbackDLQ = "helpdesk.feedback.dlq"

	// no need dlq | fire and forget
	QueuePasswordReset = "helpdesk.password-reset"
)
