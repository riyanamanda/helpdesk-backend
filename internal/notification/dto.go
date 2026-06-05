package notification

type Message struct {
	To       string
	Subject  string
	Body     string // HTML
	TextBody string // plain text fallback (required for spam filters)
}
