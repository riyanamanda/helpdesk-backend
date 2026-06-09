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
}

type notifier struct {
	ch        *amqp.Channel
	queueName string
}

func NewNotifier(ch *amqp.Channel) Notifier {
	return &notifier{ch: ch, queueName: QueueNewTicketEmail}
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

	// Publish to the default exchange ("").
	// With the default exchange, the routing key IS the queue name —
	// RabbitMQ delivers the message directly to the matching queue.
	// DeliveryMode: amqp.Persistent writes the message to disk so it
	// survives a broker restart (equivalent to asynq's default persistence).
	if err := n.ch.PublishWithContext(ctx,
		"",
		n.queueName,
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
