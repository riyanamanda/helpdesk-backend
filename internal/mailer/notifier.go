package mailer

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type Notifier interface {
	NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID)
}

type notifier struct {
	mailer   MailerService
	userRepo user.UserRepository
}

func NewNotifier(m MailerService, userRepo user.UserRepository) Notifier {
	return &notifier{mailer: m, userRepo: userRepo}
}

func (n *notifier) NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID) {
	go func() {
		ctx := context.WithoutCancel(ctx)

		submitterName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, submitterID); err == nil {
			submitterName = u.Name
		}

		adminEmails, err := n.userRepo.GetEmailsByRole(ctx, user.ADMIN)
		if err != nil {
			slog.ErrorContext(ctx, "mailer: failed to get admin emails", "error", err)
			return
		}

		if len(adminEmails) == 0 {
			return
		}

		msg := NewTicketMessage(ticketID, title, description, submitterName, adminEmails)
		if err := n.mailer.Send(ctx, msg); err != nil {
			slog.ErrorContext(ctx, "mailer: failed to send email", "error", err)
		}
	}()
}
