package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

var feedbackTypeLabels = map[string]string{
	"FEATURE_REQUEST": "Permintaan Fitur",
	"IMPROVEMENT":     "Peningkatan",
	"BUG_REPORT":      "Laporan Bug",
}

func (c *Consumer) feedbackURL() string {
	return c.appURL + "/feedback"
}

func (c *Consumer) handleFeedbackCreated(ctx context.Context, body []byte) error {
	var ev event.FeedbackCreatedEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return err
	}

	submitterName := "Pengguna"
	var submitterEmail string
	if submitter, err := c.userRepo.GetByID(ctx, ev.SubmittedByID); err == nil {
		submitterName = submitter.Name
		submitterEmail = submitter.Email
	}

	emails, err := c.userRepo.GetEmailsByRoles(ctx, rbac.ADMIN, rbac.SUPERADMIN)
	if err != nil {
		return fmt.Errorf("get admin emails: %w", err)
	}

	typeLabel := feedbackTypeLabels[ev.Type]
	if typeLabel == "" {
		typeLabel = ev.Type
	}

	data := mailer.FeedbackCreatedEmailData{
		FeedbackID:  ev.FeedbackID,
		Title:       ev.Title,
		Type:        typeLabel,
		Description: ev.Description,
		SubmittedBy: submitterName,
		FeedbackURL: c.feedbackURL(),
	}

	for i, email := range emails {
		if strings.EqualFold(email, submitterEmail) {
			continue
		}
		if err := c.mailer.SendFeedbackCreatedEmail(ctx, email, data); err != nil {
			return fmt.Errorf("send feedback email %d of %d: %w", i+1, len(emails), err)
		}
	}

	return nil
}

func (c *Consumer) handleFeedbackStatusChanged(ctx context.Context, body []byte) error {
	var ev event.FeedbackStatusChangedEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return err
	}

	if ev.Status != "DELIVERED" {
		return nil
	}

	submitter, err := c.recipient(ctx, ev.SubmitterID)
	if err != nil {
		return err
	}
	if submitter == nil {
		slog.InfoContext(ctx, "skip feedback delivered email", "feedback_id", ev.FeedbackID, "reason", "submitter missing or inactive")
		return nil
	}

	return c.mailer.SendFeedbackDeliveredEmail(ctx, submitter.Email, mailer.FeedbackDeliveredEmailData{
		FeedbackID:    ev.FeedbackID,
		Title:         ev.Title,
		SubmitterName: submitter.Name,
		ReviewedBy:    c.displayName(ctx, ev.ReviewedByID),
		FeedbackURL:   c.feedbackURL(),
	})
}
