package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type Worker struct {
	mailerSvc MailerService
	userRepo  user.UserRepository
}

func NewWorker(mailerSvc MailerService, userRepo user.UserRepository) *Worker {
	return &Worker{mailerSvc: mailerSvc, userRepo: userRepo}
}

func (w *Worker) HandleNewTicketEmail(ctx context.Context, d amqp.Delivery) error {
	var p newTicketPayload
	if err := json.Unmarshal(d.Body, &p); err != nil {
		_ = d.Nack(false, false)
		return fmt.Errorf("mailer: unmarshal payload: %w", err)
	}

	submitterID, err := uuid.Parse(p.SubmitterID)
	if err != nil {
		_ = d.Nack(false, false)
		return fmt.Errorf("mailer: invalid submitter_id: %w", err)
	}

	submitterName := "Unknown"
	if u, err := w.userRepo.GetByID(ctx, submitterID); err == nil {
		submitterName = u.Name
	}

	adminEmails, err := w.userRepo.GetEmailsByRole(ctx, rbac.ADMIN)
	if err != nil {
		_ = d.Nack(false, true)
		return fmt.Errorf("mailer: get admin emails: %w", err)
	}
	if len(adminEmails) == 0 {
		_ = d.Ack(false)
		return nil
	}

	msg := NewTicketMessage(p.TicketID, p.Title, p.Description, submitterName, adminEmails)
	if err := w.mailerSvc.Send(ctx, msg); err != nil {
		slog.ErrorContext(ctx, "mailer: failed to send email", "error", err)
		_ = d.Nack(false, true)
		return fmt.Errorf("mailer: send: %w", err)
	}

	_ = d.Ack(false)
	return nil
}
