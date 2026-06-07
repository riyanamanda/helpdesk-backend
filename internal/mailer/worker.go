package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type Worker struct {
	mailerSvc MailerService
	userRepo  user.UserRepository
}

func NewWorker(mailerSvc MailerService, userRepo user.UserRepository) *Worker {
	return &Worker{mailerSvc: mailerSvc, userRepo: userRepo}
}

func (w *Worker) HandleNewTicketEmail(ctx context.Context, t *asynq.Task) error {
	var p newTicketPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("mailer: unmarshal payload: %w", err)
	}

	submitterID, err := uuid.Parse(p.SubmitterID)
	if err != nil {
		return fmt.Errorf("mailer: invalid submitter_id: %w", err)
	}

	submitterName := "Unknown"
	if u, err := w.userRepo.GetByID(ctx, submitterID); err == nil {
		submitterName = u.Name
	}

	adminEmails, err := w.userRepo.GetEmailsByRole(ctx, user.ADMIN)
	if err != nil {
		return fmt.Errorf("mailer: get admin emails: %w", err)
	}
	if len(adminEmails) == 0 {
		return nil
	}

	msg := NewTicketMessage(p.TicketID, p.Title, p.Description, submitterName, adminEmails)
	if err := w.mailerSvc.Send(ctx, msg); err != nil {
		slog.ErrorContext(ctx, "mailer: failed to send email", "error", err)
		return fmt.Errorf("mailer: send: %w", err)
	}

	return nil
}
