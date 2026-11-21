package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type handlerFunc func(ctx context.Context, body []byte) error

type Consumer struct {
	rabbitmq *rabbitmq.Client
	mailer   *mailer.Mailer
	userRepo user.Repository
	queue    string
	appURL   string
	handlers map[string]handlerFunc
}

func NewConsumer(rabbitmq *rabbitmq.Client, mailer *mailer.Mailer, userRepo user.Repository, queue string, appURL string) *Consumer {
	c := &Consumer{
		rabbitmq: rabbitmq,
		mailer:   mailer,
		userRepo: userRepo,
		queue:    queue,
		appURL:   appURL,
	}

	c.handlers = map[string]handlerFunc{
		event.UserCreated:            c.handleUserCreated,
		event.TicketCreated:          c.handleTicketCreated,
		event.TicketAssigned:         c.handleTicketAssigned,
		event.TicketResolved:         c.handleTicketResolved,
		event.PasswordResetRequested: c.handlePasswordResetRequested,
		event.PasswordChanged:        c.handlePasswordChanged,
		event.FeedbackCreated:        c.handleFeedbackCreated,
		event.FeedbackStatusChanged:  c.handleFeedbackStatusChanged,
	}

	return c
}

func (c *Consumer) Run(ctx context.Context) error {
	messages, err := c.rabbitmq.Consume(c.queue)
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
				return errors.New("delivery channel closed")
			}

			handle, known := c.handlers[msg.RoutingKey]
			if !known {
				slog.Error("unknown event", "routing_key", msg.RoutingKey)
				if err := msg.Nack(false, false); err != nil {
					slog.Error("nack message", "routing_key", msg.RoutingKey, "error", err)
				}
				continue
			}

			if err := handle(ctx, msg.Body); err != nil {
				slog.Error("handle message", "routing_key", msg.RoutingKey, "error", err)
				if err := msg.Nack(false, false); err != nil {
					slog.Error("nack message", "routing_key", msg.RoutingKey, "error", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				slog.Error("ack message", "routing_key", msg.RoutingKey, "error", err)
				continue
			}

			slog.Info("message acknowledged", "routing_key", msg.RoutingKey)
		}
	}
}

func (c *Consumer) recipient(ctx context.Context, userID uuid.UUID) (*user.UserProjection, error) {
	u, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get recipient: %w", err)
	}

	if !u.IsActive {
		return nil, nil
	}

	return u, nil
}

func (c *Consumer) displayName(ctx context.Context, userID uuid.UUID) string {
	u, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "Tim IT"
	}
	return u.Name
}
