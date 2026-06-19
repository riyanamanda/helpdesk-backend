package mailer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Notifier interface {
	NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID)
	WelcomeUserEmail(ctx context.Context, name, email, password string)
}

type notifier struct {
	ticketCh  *amqp.Channel
	welcomeCh *amqp.Channel
}

func NewNotifier(ticketCh *amqp.Channel, welcomeCh *amqp.Channel) Notifier {
	return &notifier{ticketCh: ticketCh, welcomeCh: welcomeCh}
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

	if err := n.ticketCh.PublishWithContext(ctx,
		"",
		QueueNewTicketEmail,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		},
	); err != nil {
		slog.ErrorContext(ctx, "mailer: publish message", "error", err)
	}
}

func (n *notifier) WelcomeUserEmail(ctx context.Context, name, email, password string) {
	payload, err := json.Marshal(welcomeUserPayload{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		slog.ErrorContext(ctx, "mailer: marshal welcome payload", "error", err)
		return
	}

	if err := n.welcomeCh.PublishWithContext(ctx,
		"",
		QueueWelcomeUserEmail,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		},
	); err != nil {
		slog.ErrorContext(ctx, "mailer: publish welcome message", "error", err)
	}
}
