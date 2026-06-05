package email

import "github.com/riyanamanda/helpdesk-backend/internal/platform/config"

type SMTPClient struct {
	cfg config.Email
}

func NewSMTPClient(cfg config.Email) *SMTPClient {
	return &SMTPClient{cfg: cfg}
}

func (s *SMTPClient) Config() config.Email {
	return s.cfg
}
