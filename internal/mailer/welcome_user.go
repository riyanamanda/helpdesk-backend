package mailer

import "fmt"

func NewWelcomeUserMessage(name, email, password string) (Message, error) {
	data := WelcomeUserTemplateData{
		Name:     name,
		Email:    email,
		Password: password,
	}

	htmlBody, err := renderWelcomeUserHTML(data)
	if err != nil {
		return Message{}, fmt.Errorf("render welcome user HTML: %w", err)
	}

	return Message{
		To:      email,
		Subject: "Selamat Datang di IT Helpdesk — Akun Anda Telah Aktif",
		Body:    htmlBody,
	}, nil
}
