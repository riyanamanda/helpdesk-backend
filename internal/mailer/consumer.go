package mailer

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch        *amqp.Channel
	queueName string
	worker    *Worker
}

func NewConsumer(ch *amqp.Channel, queueName string, worker *Worker) *Consumer {
	return &Consumer{ch: ch, queueName: queueName, worker: worker}
}

func (c *Consumer) Start(ctx context.Context) error {
	deliveries, err := c.ch.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	slog.Info("mailer consumer started", "queue", c.queueName)

	for d := range deliveries {
		var handleErr error
		switch c.queueName {
		case QueueNewTicketEmail:
			handleErr = c.worker.HandleNewTicketEmail(ctx, d)
		case QueueWelcomeUserEmail:
			handleErr = c.worker.HandleWelcomeUserEmail(ctx, d)
		}
		if handleErr != nil {
			slog.ErrorContext(ctx, "mailer: handle delivery", "queue", c.queueName, "error", handleErr)
		}
	}

	slog.Info("mailer consumer stopped", "queue", c.queueName)
	return nil
}

func (c *Consumer) Shutdown() {
	if err := c.ch.Close(); err != nil {
		slog.Error("mailer: consumer channel close", "error", err)
	}
}
