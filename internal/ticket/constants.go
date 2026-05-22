package ticket

type TicketStatus string
type TicketPriority string

type AttachmentType string

const (
	Open       TicketStatus = "OPEN"
	InProgress TicketStatus = "IN_PROGRESS"
	Resolved   TicketStatus = "RESOLVED"
	Closed     TicketStatus = "CLOSED"
)

const (
	Low    TicketPriority = "LOW"
	Medium TicketPriority = "MEDIUM"
	High   TicketPriority = "HIGH"
	Urgent TicketPriority = "URGENT"
)

const (
	Report     AttachmentType = "REPORT"
	Resolution AttachmentType = "RESOLUTION"
)

const maxTicketAttachmentSize = 2 << 20 // 2MB

var AllowedTicketAttachmentTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}
