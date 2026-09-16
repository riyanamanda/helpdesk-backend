package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Notifier interface {
	NewTicketEmail(ctx context.Context, ticketID int64, title, description string, submitterID uuid.UUID)
	WelcomeUserEmail(ctx context.Context, name, email, password string)
}

type notifier struct {
	dialer    Dialer
	mu        sync.Mutex
	conn      *amqp.Connection
	ticketCh  *amqp.Channel
	welcomeCh *amqp.Channel
}

func NewNotifier(dialer Dialer) Notifier {
	return &notifier{dialer: dialer}
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

	if err := n.publish(ctx, QueueNewTicketEmail, payload); err != nil {
		slog.ErrorContext(ctx, "mailer: publish new ticket email", "error", err)
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

	if err := n.publish(ctx, QueueWelcomeUserEmail, payload); err != nil {
		slog.ErrorContext(ctx, "mailer: publish welcome email", "error", err)
	}
}

// publish reconnects if the connection is dead and retries with backoff.
func (n *notifier) publish(ctx context.Context, queueName string, body []byte) error {
	const maxAttempts = 3

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		n.mu.Lock()

		if n.conn == nil || n.conn.IsClosed() {
			if err := n.connectLocked(); err != nil {
				n.mu.Unlock()
				lastErr = err
				if err := sleep(ctx, attempt); err != nil {
					return err
				}
				continue
			}
		}

		ch := n.welcomeCh
		if queueName == QueueNewTicketEmail {
			ch = n.ticketCh
		}

		err := ch.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})

		if err == nil {
			n.mu.Unlock()
			return nil
		}

		// Mark the connection as dead so the next attempt reconnects.
		lastErr = err
		_ = n.conn.Close()
		n.conn = nil
		n.ticketCh = nil
		n.welcomeCh = nil
		n.mu.Unlock()

		if err := sleep(ctx, attempt); err != nil {
			return err
		}
	}

	return fmt.Errorf("publish %q failed after %d attempts: %w", queueName, maxAttempts, lastErr)
}

// connectLocked (re)establishes the connection and both channels.
// The caller must hold n.mu.
func (n *notifier) connectLocked() error {
	if n.conn != nil {
		_ = n.conn.Close()
	}

	conn, err := n.dialer()
	if err != nil {
		return err
	}

	ticketCh, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if _, err := ticketCh.QueueDeclare(QueueNewTicketEmail, true, false, false, false, nil); err != nil {
		_ = ticketCh.Close()
		_ = conn.Close()
		return err
	}

	welcomeCh, err := conn.Channel()
	if err != nil {
		_ = ticketCh.Close()
		_ = conn.Close()
		return err
	}
	if _, err := welcomeCh.QueueDeclare(QueueWelcomeUserEmail, true, false, false, false, nil); err != nil {
		_ = welcomeCh.Close()
		_ = ticketCh.Close()
		_ = conn.Close()
		return err
	}

	n.conn = conn
	n.ticketCh = ticketCh
	n.welcomeCh = welcomeCh
	return nil
}

func sleep(ctx context.Context, attempt int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Duration(attempt+1) * time.Second):
		return nil
	}
}
