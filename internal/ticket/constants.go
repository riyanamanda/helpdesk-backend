package ticket

type TicketStatus string
type TicketPriority string

type AttachmentType string

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

const (
	REPORT     AttachmentType = "REPORT"
	RESOLUTION AttachmentType = "RESOLUTION"
)

const maxTicketAttachmentSize = 2 << 20 // 2MB

var AllowedTicketAttachmentTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}
