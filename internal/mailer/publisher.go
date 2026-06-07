package mailer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type Notifier interface {
	NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID)
}

type notifier struct {
	client *asynq.Client
}

func NewNotifier(client *asynq.Client) Notifier {
	return &notifier{client: client}
}

func (n *notifier) NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID) {
	payload, err := json.Marshal(newTicketPayload{
		TicketID:    ticketID,
		Title:       title,
		Description: description,
		SubmitterID: submitterID.String(),
	})
	if err != nil {
		slog.ErrorContext(ctx, "mailer: marshal payload", "error", err)
		return
	}

	task := asynq.NewTask(TaskNewTicketEmail, payload)
	if _, err := n.client.EnqueueContext(ctx, task); err != nil {
		slog.ErrorContext(ctx, "mailer: enqueue task", "error", err)
	}
}
