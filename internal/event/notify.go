package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Publisher interface {
	PublishJSON(ctx context.Context, routingKey string, v any) error
}

func NotifyPasswordChanged(ctx context.Context, p Publisher, userID uuid.UUID, byAdmin bool) {
	err := p.PublishJSON(ctx, PasswordChanged, PasswordChangedEvent{
		UserID:    userID,
		ChangedAt: time.Now(),
		ByAdmin:   byAdmin,
	})
	if err != nil {
		slog.WarnContext(ctx, "publish password changed", "user_id", userID, "error", err)
	}
}
