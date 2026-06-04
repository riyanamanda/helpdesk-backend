package ticket

type (
	TicketStatus   string
	TicketPriority string
	AttachmentType string
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
