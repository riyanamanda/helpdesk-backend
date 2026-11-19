package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
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

	return c.mailer.SendResetPasswordEmail(ctx, passwordRequestedEvent.Name, passwordRequestedEvent.Email, passwordRequestedEvent.ResetURL, passwordRequestedEvent.ExpiresIn)
}

var priorityLabels = map[string]string{
	"LOW":    "Rendah",
	"MEDIUM": "Sedang",
	"HIGH":   "Tinggi",
	"URGENT": "Mendesak",
}

var wib = time.FixedZone("WIB", 7*60*60)

func (c *Consumer) ticketURL(ticketID int64) string {
	return fmt.Sprintf("%s/ticket/%d", c.appURL, ticketID)
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

func (c *Consumer) handleTicketAssigned(ctx context.Context, body []byte) error {
	var ev event.TicketAssignedEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return err
	}

	assignee, err := c.recipient(ctx, ev.AssigneeID)
	if err != nil {
		return err
	}
	if assignee == nil {
		slog.InfoContext(ctx, "skip ticket assigned email", "ticket_id", ev.TicketID, "reason", "assignee missing or inactive")
		return nil
	}

	priority := priorityLabels[ev.Priority]
	if priority == "" {
		priority = ev.Priority
	}

	return c.mailer.SendTicketAssignedEmail(ctx, assignee.Email, mailer.TicketAssignedEmailData{
		TicketID:   ev.TicketID,
		Title:      ev.Title,
		Priority:   priority,
		AssignedBy: c.displayName(ctx, ev.AssignedByID),
		Note:       ev.Note,
		TicketURL:  c.ticketURL(ev.TicketID),
	})
}

func (c *Consumer) handleTicketResolved(ctx context.Context, body []byte) error {
	var ev event.TicketResolvedEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return err
	}

	reporter, err := c.recipient(ctx, ev.ReporterID)
	if err != nil {
		return err
	}
	if reporter == nil {
		slog.InfoContext(ctx, "skip ticket resolved email", "ticket_id", ev.TicketID, "reason", "reporter missing or inactive")
		return nil
	}

	return c.mailer.SendTicketResolvedEmail(ctx, reporter.Email, mailer.TicketResolvedEmailData{
		TicketID:     ev.TicketID,
		Title:        ev.Title,
		ReporterName: reporter.Name,
		ResolvedBy:   c.displayName(ctx, ev.ResolvedByID),
		Resolution:   ev.Resolution,
		TicketURL:    c.ticketURL(ev.TicketID),
	})
}

func (c *Consumer) handlePasswordChanged(ctx context.Context, body []byte) error {
	var ev event.PasswordChangedEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return err
	}

	u, err := c.recipient(ctx, ev.UserID)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}

	return c.mailer.SendPasswordChangedEmail(ctx, u.Email, mailer.PasswordChangedEmailData{
		Name:      u.Name,
		ChangedAt: ev.ChangedAt.In(wib).Format("02/01/2006 15.04") + " WIB",
		ByAdmin:   ev.ByAdmin,
		LoginURL:  c.appURL + "/login",
	})
}
