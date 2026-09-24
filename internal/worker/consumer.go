package worker

import (
	"encoding/json"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
)

type Consumer struct {
	rabbitmq *rabbitmq.Client
}

func NewConsumer(rabbitmq *rabbitmq.Client) *Consumer {
	return &Consumer{
		rabbitmq: rabbitmq,
	}
}

func (c *Consumer) Run() error {
	messages, err := c.rabbitmq.Consume("helpdesk.email")
	if err != nil {
		return err
	}

	for msg := range messages {
		var event UserCreatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			slog.Error("decode user.created event failed", "error", err)

			msg.Nack(false, false)
			continue
		}

		slog.Info("user.created event received", "name", event.Name, "email", event.Email)

		// TODO: send welcome email

		if err := msg.Ack(false); err != nil {
			slog.Error("message ack failed", "error", err)
			continue
		}

		slog.Info("message acknoledged")
	}

	return nil
}
