package worker

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
)

type Consumer struct {
	rabbitmq *rabbitmq.Client
	mailer   *mailer.Mailer
}

func NewConsumer(rabbitmq *rabbitmq.Client, mailer *mailer.Mailer) *Consumer {
	return &Consumer{
		rabbitmq: rabbitmq,
		mailer:   mailer,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	messages, err := c.rabbitmq.Consume("helpdesk.email")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("consumer stopped")
			return nil

		case msg, ok := <-messages:
			if !ok {
				slog.Error("rabbitmq message closed")
				return nil
			}

			var event UserCreatedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				slog.Error("decode user.created event failed", "error", err)

				msg.Nack(false, false)
				continue
			}

			if err := c.mailer.SendWelcomeEmail(ctx, event.Name, event.Email); err != nil {
				slog.Error("send welcome email failed", "error", err)
				msg.Nack(false, false)
				continue
			}

			if err := msg.Ack(false); err != nil {
				slog.Error("message ack failed", "error", err)
				continue
			}

			slog.Info("message acknoledged")
		}
	}
}
