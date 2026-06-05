package notification

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type NotificationService interface {
	NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID)
}

type service struct {
	mailer   mailer.MailerService
	userRepo user.UserRepository
}

func NewNotificationService(m mailer.MailerService, userRepo user.UserRepository) NotificationService {
	return &service{
		mailer:   m,
		userRepo: userRepo,
	}
}

func (s *service) NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID) {
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

		if len(adminEmails) == 0 {
			return
		}

		msg := mailer.NewTicketMessage(ticketID, title, description, submitterName, adminEmails)
		if err := s.mailer.Send(ctx, msg); err != nil {
			slog.ErrorContext(ctx, "notification: failed to send email", "error", err)
		}
	}()
}
