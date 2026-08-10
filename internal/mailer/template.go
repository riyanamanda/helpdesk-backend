package mailer

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/ticket/created.html
//go:embed templates/user/welcome.html
var templateFS embed.FS

type TicketCreatedTemplateData struct {
	TicketID    int64
	Title       string
	Description string
	SubmittedBy string
}

type WelcomeUserTemplateData struct {
	Name     string
	Email    string
	Password string
}

func renderTicketCreatedHTML(data TicketCreatedTemplateData) (string, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/ticket/created.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
func renderWelcomeUserHTML(data WelcomeUserTemplateData) (string, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/user/welcome.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
