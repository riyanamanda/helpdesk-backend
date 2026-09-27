package worker

import (
	"context"
	"encoding/json"
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
	userRepo user.UserRepository
	queue    string
}

func NewConsumer(rabbitmq *rabbitmq.Client, mailer *mailer.Mailer, userRepo user.UserRepository, queue string) *Consumer {
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
				slog.Error("rabbitmq message closed")
				return nil
			}

			switch msg.RoutingKey {
			case event.UserCreated:
				if err := c.handleUserCreated(ctx, msg.Body); err != nil {
					slog.Error("handle user.created failed", "error", err)
					if err := msg.Nack(false, false); err != nil {
						slog.Error("message nack failed", "error", err)
					}
					continue
				}

			case event.TicketCreated:
				if err := c.handleTicketCreated(ctx, msg.Body); err != nil {
					slog.Error("handle ticket.created failed", "error", err)
					if err := msg.Nack(false, false); err != nil {
						slog.Error("message nack failed", "error", err)
					}
					continue
				}

			default:
				slog.Error("unknown event", "routing_key", msg.RoutingKey)
				if err := msg.Nack(false, false); err != nil {
					slog.Error("message nack failed, unknown event", "error", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				slog.Error("message ack failed", "error", err)
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

	return c.mailer.SendWelcomeEmail(ctx, createdEvent.Name, createdEvent.Email)
}

func (c *Consumer) handleTicketCreated(ctx context.Context, body []byte) error {
	var createdEvent event.TicketCreatedEvent

	if err := json.Unmarshal(body, &createdEvent); err != nil {
		return err
	}

	emails, err := c.userRepo.GetEmailsByRoles(ctx, rbac.ADMIN, rbac.SUPERADMIN)
	if err != nil {
		return fmt.Errorf("get admin emails, %w", err)
	}

	for _, email := range emails {
		if err := c.mailer.SendNewTicketEmail(ctx, email, createdEvent.TicketID, createdEvent.SubmittedBy, createdEvent.Title, createdEvent.Description); err != nil {
			return fmt.Errorf("send ticket email to %s: %w", email, err)
		}
	}

	return nil
}
