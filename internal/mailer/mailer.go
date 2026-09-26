package mailer

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
)

type Mailer struct {
	config config.Email
}

type WelcomeEmailData struct {
	Name  string
	Email string
}

type TicketEmailData struct {
	TicketID    int64
	SubmittedBy string
	Title       string
	Description string
}

//go:embed templates/user/welcome.html
var welcomeTemplate string

//go:embed templates/ticket/created.html
var ticketCreatedTemplate string

func NewMailer(config config.Email) *Mailer {
	return &Mailer{
		config: config,
	}
}

func (m *Mailer) SendWelcomeEmail(ctx context.Context, name string, email string) error {
	tmpl, err := template.New("welcome.html").Parse(welcomeTemplate)
	if err != nil {
		return fmt.Errorf("parse welcome email template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, WelcomeEmailData{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return fmt.Errorf("execute welcome email template: %w", err)
	}

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	message := []byte(
		"From: " + m.config.From + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Selamat Datang di IT Helpdesk\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String(),
	)

	return smtp.SendMail(
		m.config.Host+":"+m.config.Port,
		auth,
		m.config.From,
		[]string{email},
		message,
	)
}

func (m *Mailer) SendNewTicketEmail(ctx context.Context, email string, ticketID int64, submittedBy string, title string, description string) error {
	tmpl, err := template.New("created.html").Parse(ticketCreatedTemplate)
	if err != nil {
		return fmt.Errorf("parse ticket email template: %w", err)
	}

	var body bytes.Buffer

	err = tmpl.Execute(&body, TicketEmailData{
		TicketID:    ticketID,
		SubmittedBy: submittedBy,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return fmt.Errorf("execute ticket email template: %w", err)
	}

	auth := smtp.PlainAuth(
		"",
		m.config.Username,
		m.config.Password,
		m.config.Host,
	)

	message := []byte(
		"From: " + m.config.From + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Tiket Dukungan Baru #" + strconv.FormatInt(ticketID, 10) + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String(),
	)

	return smtp.SendMail(
		m.config.Host+":"+m.config.Port,
		auth,
		m.config.From,
		[]string{email},
		message,
	)
}
