package notification

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/email"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type NotificationService interface {
	NewTicket(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID)
}

type service struct {
	client   *email.SMTPClient
	userRepo user.UserRepository
}

func NewNotificationService(client *email.SMTPClient, userRepo user.UserRepository) NotificationService {
	return &service{
		client:   client,
		userRepo: userRepo,
	}
}

func (s *service) NewTicket(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID) {
	go func() {
		ctx := context.WithoutCancel(ctx)

		submitterName := "Unknown"
		if u, err := s.userRepo.GetByID(ctx, submitterID); err == nil {
			submitterName = u.Name
		}

		adminEmails, err := s.userRepo.GetEmailsByRole(ctx, user.ADMIN)
		if err != nil {
			slog.ErrorContext(ctx, "notification: failed to get admin emails", "error", err)
			return
		}

		subject := fmt.Sprintf("New Ticket #%d: %s", ticketID, title)
		msg := Message{
			Subject:  subject,
			Body:     newTicketHTMLBody(ticketID, title, description, submitterName),
			TextBody: newTicketTextBody(ticketID, title, description, submitterName),
		}

		for _, to := range adminEmails {
			msg.To = to
			if err := s.send(ctx, msg); err != nil {
				slog.ErrorContext(ctx, "notification: failed to send email", "to", to, "error", err)
			}
		}
	}()
}

func (s *service) send(_ context.Context, msg Message) error {
	cfg := s.client.Config()
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

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
	header := fmt.Sprintf(
		"From: Helpdesk <%s>\r\n"+
			"Reply-To: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Date: %s\r\n"+
			"Message-ID: <%d.%d@helpdesk>\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/alternative; boundary=%q\r\n\r\n",
		cfg.From, cfg.From, msg.To, msg.Subject,
		now.Format(time.RFC1123Z),
		now.UnixNano(), now.Unix(),
		mw.Boundary(),
	)

	return smtp.SendMail(addr, auth, cfg.Username, []string{msg.To}, []byte(header+body.String()))
}
