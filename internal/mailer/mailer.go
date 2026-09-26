package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"html/template"
	"net"
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

func (m *Mailer) send(ctx context.Context, email string, message []byte) error {
	addr := net.JoinHostPort(m.config.Host, m.config.Port)

	if m.config.UseSSL {
		tlsConfig := &tls.Config{
			ServerName: m.config.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("connect SMTP with TLS: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.config.Host)
		if err != nil {
			return fmt.Errorf("create SMTP client: %w", err)
		}
		defer client.Close()

		auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication: %w", err)
		}

		if err := client.Mail(m.config.From); err != nil {
			return fmt.Errorf("SMTP MAIL FROM: %w", err)
		}

		if err := client.Rcpt(email); err != nil {
			return fmt.Errorf("SMTP RCPT TO: %w", err)
		}

		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("SMTP DATA: %w", err)
		}

		if _, err := writer.Write(message); err != nil {
			writer.Close()
			return fmt.Errorf("write SMTP message: %w", err)
		}

		if err := writer.Close(); err != nil {
			return fmt.Errorf("close SMTP message: %w", err)
		}

		return client.Quit()
	}

	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

	return smtp.SendMail(
		addr,
		auth,
		m.config.From,
		[]string{email},
		message,
	)
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

	message := []byte(
		"From: IT Helpdesk <" + m.config.From + ">\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Welcome to IT Helpdesk\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String(),
	)

	return m.send(ctx, email, message)
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

	message := []byte(
		"From: IT Helpdesk <" + m.config.From + ">\r\n" +
			"To: " + email + "\r\n" +
			"Subject: [IT Helpdesk] New Ticket #" + strconv.FormatInt(ticketID, 10) + " — Need Review ASAP\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String(),
	)

	return m.send(ctx, email, message)
}
