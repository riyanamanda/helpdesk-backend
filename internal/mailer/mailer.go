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
	Name           string
	Email          string
	SetPasswordURL string
	ExpiresIn      string
}

type TicketEmailData struct {
	TicketID    int64
	SubmittedBy string
	Title       string
	Description string
}

type ResetPasswordEmailData struct {
	Name      string
	ResetURL  string
	ExpiresIn string
}

type TicketAssignedEmailData struct {
	TicketID   int64
	Title      string
	Priority   string
	AssignedBy string
	Note       string
	TicketURL  string
}

type TicketResolvedEmailData struct {
	TicketID     int64
	Title        string
	ReporterName string
	ResolvedBy   string
	Resolution   string
	TicketURL    string
}

type PasswordChangedEmailData struct {
	Name      string
	ChangedAt string
	ByAdmin   bool
	LoginURL  string
}

type FeedbackCreatedEmailData struct {
	FeedbackID  int64
	Title       string
	Type        string
	Description string
	SubmittedBy string
	FeedbackURL string
}

type FeedbackDeliveredEmailData struct {
	FeedbackID    int64
	Title         string
	SubmitterName string
	ReviewedBy    string
	FeedbackURL   string
}

const defaultResetExpiry = "15 menit"

//go:embed templates/user/welcome.html
var welcomeTemplate string

//go:embed templates/user/reset-password.html
var resetPasswordTemplate string

//go:embed templates/user/password-changed.html
var passwordChangedTemplate string

//go:embed templates/ticket/created.html
var ticketCreatedTemplate string

//go:embed templates/ticket/assigned.html
var ticketAssignedTemplate string

//go:embed templates/ticket/resolved.html
var ticketResolvedTemplate string

//go:embed templates/feedback/created.html
var feedbackCreatedTemplate string

//go:embed templates/feedback/delivered.html
var feedbackDeliveredTemplate string

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

func (m *Mailer) SendWelcomeEmail(ctx context.Context, name string, email string, setPasswordURL string, expiresIn string) error {
	tmpl, err := template.New("welcome.html").Parse(welcomeTemplate)
	if err != nil {
		return fmt.Errorf("parse welcome email template: %w", err)
	}

	var body bytes.Buffer

	err = tmpl.Execute(&body, WelcomeEmailData{
		Name:           name,
		Email:          email,
		SetPasswordURL: setPasswordURL,
		ExpiresIn:      expiresIn,
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

func (m *Mailer) SendResetPasswordEmail(ctx context.Context, name string, email string, resetUrl string, expiresIn string) error {
	tmpl, err := template.New("reset-password.html").Parse(resetPasswordTemplate)
	if err != nil {
		return fmt.Errorf("parse reset password email template: %w", err)
	}

	if expiresIn == "" {
		expiresIn = defaultResetExpiry
	}

	var body bytes.Buffer

	err = tmpl.Execute(&body, ResetPasswordEmailData{
		Name:      name,
		ResetURL:  resetUrl,
		ExpiresIn: expiresIn,
	})
	if err != nil {
		return fmt.Errorf("execute reset password email template: %w", err)
	}

	message := []byte(
		"From: IT Helpdesk <" + m.config.From + ">\r\n" +
			"To: " + email + "\r\n" +
			"Subject: [IT Helpdesk] Reset Password\r\n" +
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

func (m *Mailer) SendTicketAssignedEmail(ctx context.Context, email string, data TicketAssignedEmailData) error {
	subject := "[IT Helpdesk] Tiket #" + strconv.FormatInt(data.TicketID, 10) + " ditugaskan kepada Anda"
	return m.sendTemplate(ctx, email, subject, "assigned.html", ticketAssignedTemplate, data)
}

func (m *Mailer) SendTicketResolvedEmail(ctx context.Context, email string, data TicketResolvedEmailData) error {
	subject := "[IT Helpdesk] Tiket #" + strconv.FormatInt(data.TicketID, 10) + " telah diselesaikan"
	return m.sendTemplate(ctx, email, subject, "resolved.html", ticketResolvedTemplate, data)
}

func (m *Mailer) SendFeedbackCreatedEmail(ctx context.Context, email string, data FeedbackCreatedEmailData) error {
	subject := "[IT Helpdesk] Feedback baru #" + strconv.FormatInt(data.FeedbackID, 10)
	return m.sendTemplate(ctx, email, subject, "feedback-created.html", feedbackCreatedTemplate, data)
}

func (m *Mailer) SendFeedbackDeliveredEmail(ctx context.Context, email string, data FeedbackDeliveredEmailData) error {
	subject := "[IT Helpdesk] Feedback #" + strconv.FormatInt(data.FeedbackID, 10) + " telah dirilis"
	return m.sendTemplate(ctx, email, subject, "feedback-delivered.html", feedbackDeliveredTemplate, data)
}

func (m *Mailer) SendPasswordChangedEmail(ctx context.Context, email string, data PasswordChangedEmailData) error {
	return m.sendTemplate(ctx, email, "[IT Helpdesk] Kata sandi Anda telah diubah", "password-changed.html", passwordChangedTemplate, data)
}

func (m *Mailer) sendTemplate(ctx context.Context, email string, subject string, name string, text string, data any) error {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return fmt.Errorf("parse %s: %w", name, err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("execute %s: %w", name, err)
	}

	message := []byte(
		"From: IT Helpdesk <" + m.config.From + ">\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String(),
	)

	return m.send(ctx, email, message)
}
