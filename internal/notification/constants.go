package notification

type NotificationType string
type NotificationReferenceType string

const (
	NewTicket             NotificationType = "NEW_TICKET"
	TicketAssigned        NotificationType = "TICKET_ASSIGNED"
	TicketInProgress      NotificationType = "TICKET_IN_PROGRESS"
	TicketClosed          NotificationType = "TICKET_CLOSED"
	FeedbackStatusUpdated NotificationType = "FEEDBACK_STATUS_UPDATED"

	TicketReferenceType   NotificationReferenceType = "TICKET"
	FeedbackReferenceType NotificationReferenceType = "FEEDBACK"
)
