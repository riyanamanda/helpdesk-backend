package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/email"
)

type MailerService interface {
	Send(ctx context.Context, msg Message) error
}

type service struct {
	client *email.SMTPClient
}

func NewMailerService(client *email.SMTPClient) MailerService {
	return &service{
		client: client,
	}
}

func (s *service) Send(_ context.Context, msg Message) error {
	cfg := s.client.Config()
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	raw := s.buildRaw(cfg.From, msg)
	recipients := append([]string{msg.To}, msg.CC...)

	if cfg.UseSSL {
		// port 465: implicit TLS — must dial TLS first, then hand to smtp.NewClient
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: cfg.Host})
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return err
		}
		defer client.Close()

		if err := client.Auth(auth); err != nil {
			return err
		}
		if err := client.Mail(cfg.From); err != nil {
			return err
		}
		for _, rcpt := range recipients {
			if err := client.Rcpt(rcpt); err != nil {
				return err
			}
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write(raw); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		return client.Quit()
	}

	// port 587: STARTTLS — smtp.SendMail handles the upgrade automatically
	return smtp.SendMail(addr, auth, cfg.From, recipients, raw)
}

func (s *service) buildRaw(from string, msg Message) []byte {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	// plain text part — spam filters expect this
	textPart, _ := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qw := quotedprintable.NewWriter(textPart)
	_, _ = qw.Write([]byte(msg.TextBody))
	_ = qw.Close()

	// HTML part
	htmlPart, _ := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qw = quotedprintable.NewWriter(htmlPart)
	_, _ = qw.Write([]byte(msg.Body))
	_ = qw.Close()

	_ = mw.Close()

	now := time.Now()
	var hdr strings.Builder
	fmt.Fprintf(&hdr, "From: IT Helpdesk <%s>\r\n", from)
	fmt.Fprintf(&hdr, "Reply-To: %s\r\n", from)
	fmt.Fprintf(&hdr, "To: %s\r\n", msg.To)
	if len(msg.CC) > 0 {
		fmt.Fprintf(&hdr, "Cc: %s\r\n", strings.Join(msg.CC, ", "))
	}
	fmt.Fprintf(&hdr, "Subject: %s\r\n", msg.Subject)
	fmt.Fprintf(&hdr, "Date: %s\r\n", now.Format(time.RFC1123Z))
	fmt.Fprintf(&hdr, "Message-ID: <%d.%d@rs-erba.go.id>\r\n", now.UnixNano(), now.Unix())
	fmt.Fprintf(&hdr, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&hdr, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", mw.Boundary())

	return []byte(hdr.String() + body.String())
}
