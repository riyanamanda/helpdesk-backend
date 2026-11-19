package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type Consumer struct {
	rabbitmq *rabbitmq.Client
	mailer   *mailer.Mailer
	userRepo user.Repository
	queue    string
}

func NewConsumer(rabbitmq *rabbitmq.Client, mailer *mailer.Mailer, userRepo user.Repository, queue string) *Consumer {
	return &Consumer{
		rabbitmq: rabbitmq,
		mailer:   mailer,
		userRepo: userRepo,
		queue:    queue,
	}
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

			switch msg.RoutingKey {
			case event.UserCreated:
				if err := c.handleUserCreated(ctx, msg.Body); err != nil {
					slog.Error("handle message", "routing_key", msg.RoutingKey, "error", err)
					if err := msg.Nack(false, false); err != nil {
						slog.Error("nack message", "routing_key", msg.RoutingKey, "error", err)
					}
					continue
				}

			case event.TicketCreated:
				if err := c.handleTicketCreated(ctx, msg.Body); err != nil {
					slog.Error("handle message", "routing_key", msg.RoutingKey, "error", err)
					if err := msg.Nack(false, false); err != nil {
						slog.Error("nack message", "routing_key", msg.RoutingKey, "error", err)
					}
					continue
				}

			case event.PasswordResetRequested:
				if err := c.handlePasswordResetRequested(ctx, msg.Body); err != nil {
					slog.Error("handle message", "routing_key", msg.RoutingKey, "error", err)
					if err := msg.Nack(false, false); err != nil {
						slog.Error("nack message", "routing_key", msg.RoutingKey, "error", err)
					}
					continue
				}

			default:
				slog.Error("unknown event", "routing_key", msg.RoutingKey)
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

func (c *Consumer) handleUserCreated(ctx context.Context, body []byte) error {
	var createdEvent event.UserCreatedEvent

	if err := json.Unmarshal(body, &createdEvent); err != nil {
		return err
	}

	return c.mailer.SendWelcomeEmail(ctx, createdEvent.Name, createdEvent.Email, createdEvent.SetPasswordURL, createdEvent.ExpiresIn)
}

func (c *Consumer) handleTicketCreated(ctx context.Context, body []byte) error {
	var createdEvent event.TicketCreatedEvent

	if err := json.Unmarshal(body, &createdEvent); err != nil {
		return err
	}

	emails, err := c.userRepo.GetEmailsByRoles(ctx, rbac.ADMIN, rbac.SUPERADMIN)
	if err != nil {
		return fmt.Errorf("get admin emails: %w", err)
	}

	for i, email := range emails {
		if err := c.mailer.SendNewTicketEmail(ctx, email, createdEvent.TicketID, createdEvent.SubmittedBy, createdEvent.Title, createdEvent.Description); err != nil {
			return fmt.Errorf("send ticket email %d of %d: %w", i+1, len(emails), err)
		}
	}

	return nil
}

func (c *Consumer) handlePasswordResetRequested(ctx context.Context, body []byte) error {
	var passwordRequestedEvent event.PasswordResetRequestedEvent

	if err := json.Unmarshal(body, &passwordRequestedEvent); err != nil {
		return err
	}

	return c.mailer.SendResetPasswordEmail(ctx, passwordRequestedEvent.Name, passwordRequestedEvent.Email, passwordRequestedEvent.ResetURL, passwordRequestedEvent.ExpiresIn, passwordRequestedEvent.Forced)
}
