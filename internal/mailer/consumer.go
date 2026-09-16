package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Dialer opens a new RabbitMQ connection.
type Dialer func() (*amqp.Connection, error)

// Consumer consumes from a single queue and automatically reconnects whenever
// the underlying connection drops (e.g. RabbitMQ restarts).
type Consumer struct {
	dialer     Dialer
	queueName  string
	worker     *Worker
	retryDelay time.Duration
}

func NewConsumer(dialer Dialer, queueName string, worker *Worker) *Consumer {
	return &Consumer{
		dialer:     dialer,
		queueName:  queueName,
		worker:     worker,
		retryDelay: 3 * time.Second,
	}
}

// Run consumes until ctx is cancelled, reconnecting after any connection loss.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		if err := c.consumeOnce(ctx); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.ErrorContext(ctx, "mailer consumer stopped, reconnecting",
				"queue", c.queueName, "error", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(c.retryDelay):
		}
	}
}

func (c *Consumer) consumeOnce(ctx context.Context) error {
	conn, err := c.dialer()
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(c.queueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare queue %q: %w", c.queueName, err)
	}

	// Prefetch 1 message at a time (fair dispatch with manual ack).
	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("qos %q: %w", c.queueName, err)
	}

	deliveries, err := ch.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume %q: %w", c.queueName, err)
	}

	slog.Info("mailer consumer started", "queue", c.queueName)

	// connClose is signalled when the broker closes the connection.
	connClose := conn.NotifyClose(make(chan *amqp.Error, 1))

	for {
		select {
		case d, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}
			if err := c.handle(ctx, d); err != nil {
				slog.ErrorContext(ctx, "mailer: handle delivery",
					"queue", c.queueName, "error", err)
			}
		case <-connClose:
			return fmt.Errorf("connection closed")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Consumer) handle(ctx context.Context, d amqp.Delivery) error {
	switch c.queueName {
	case QueueNewTicketEmail:
		return c.worker.HandleNewTicketEmail(ctx, d)
	case QueueWelcomeUserEmail:
		return c.worker.HandleWelcomeUserEmail(ctx, d)
	default:
		_ = d.Nack(false, false)
		return fmt.Errorf("unknown queue %q", c.queueName)
	}
}
