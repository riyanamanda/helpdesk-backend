package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

var priorityLabels = map[string]string{
	"LOW":    "Rendah",
	"MEDIUM": "Sedang",
	"HIGH":   "Tinggi",
	"URGENT": "Mendesak",
}

func (c *Consumer) ticketURL(ticketID int64) string {
	return fmt.Sprintf("%s/ticket/%d", c.appURL, ticketID)
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
