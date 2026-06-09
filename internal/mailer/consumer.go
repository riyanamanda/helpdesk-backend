package mailer

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch     *amqp.Channel
	worker *Worker
}

func NewConsumer(ch *amqp.Channel, worker *Worker) *Consumer {
	return &Consumer{ch: ch, worker: worker}
}

func (c *Consumer) Start(ctx context.Context) error {
	deliveries, err := c.ch.Consume(QueueNewTicketEmail, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	slog.Info("mailer consumer started", "queue", QueueNewTicketEmail)

	for d := range deliveries {
		if err := c.worker.HandleNewTicketEmail(ctx, d); err != nil {
			slog.ErrorContext(ctx, "mailer: handle delivery", "error", err)
		}
	}

	slog.Info("mailer consumer stopped", "queue", QueueNewTicketEmail)
	return nil
}

func (c *Consumer) Shutdown() {
	if err := c.ch.Close(); err != nil {
		slog.Error("mailer: consumer channel close", "error", err)
	}
}
