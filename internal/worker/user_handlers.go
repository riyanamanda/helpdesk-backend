package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
)

var wib = time.FixedZone("WIB", 7*60*60)

func (c *Consumer) handleUserCreated(ctx context.Context, body []byte) error {
	var createdEvent event.UserCreatedEvent

	if err := json.Unmarshal(body, &createdEvent); err != nil {
		return err
	}

	return c.mailer.SendWelcomeEmail(ctx, createdEvent.Name, createdEvent.Email, createdEvent.SetPasswordURL, createdEvent.ExpiresIn)
}

func (c *Consumer) handlePasswordResetRequested(ctx context.Context, body []byte) error {
	var passwordRequestedEvent event.PasswordResetRequestedEvent

	if err := json.Unmarshal(body, &passwordRequestedEvent); err != nil {
		return err
	}

	return c.mailer.SendResetPasswordEmail(ctx, passwordRequestedEvent.Name, passwordRequestedEvent.Email, passwordRequestedEvent.ResetURL, passwordRequestedEvent.ExpiresIn)
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
