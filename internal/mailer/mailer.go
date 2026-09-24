package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
)

type Mailer struct {
	config config.Email
}

type WelcomeEmailData struct {
	Name  string
	Email string
}

func NewMailer(config config.Email) *Mailer {
	return &Mailer{
		config: config,
	}
}

func (m *Mailer) SendWelcomeEmail(ctx context.Context, name string, email string) error {
	tmpl, err := template.ParseFiles("internal/mailer/templates/user/welcome.html")
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
