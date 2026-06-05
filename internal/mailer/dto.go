package mailer

type Message struct {
	To       string
	CC       []string
	Subject  string
	Body     string
	TextBody string
}
