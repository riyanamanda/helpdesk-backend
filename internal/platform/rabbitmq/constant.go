package rabbitmq

const (
	// exchange
	ExchangeEvent      = "helpdesk.events"
	ExchangeWelcomeDLX = "helpdesk.welcome.dlx"
	ExchangeTicketDLX  = "helpdesk.ticket.dlx"

	// queue
	QueueWelcome    = "helpdesk.welcome"
	QueueWelcomeDLQ = "helpdesk.welcome.dlq"
	QueueTicket     = "helpdesk.ticket"
	QueueTicketDLQ  = "helpdesk.ticket.dlq"
)
