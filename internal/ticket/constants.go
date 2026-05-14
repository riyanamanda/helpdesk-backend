package ticket

type TicketStatus string
type TicketPriority string

const (
	OPEN        TicketStatus = "OPEN"
	IN_PROGRESS TicketStatus = "IN_PROGRESS"
	RESOLVED    TicketStatus = "RESOLVED"
	CLOSED      TicketStatus = "CLOSED"
)

const (
	LOW    TicketPriority = "LOW"
	MEDIUM TicketPriority = "MEDIUM"
	HIGH   TicketPriority = "HIGH"
	URGENT TicketPriority = "URGENT"
)
