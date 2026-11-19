package presence

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	onlineTTL    = 5 * time.Minute
	syncInterval = time.Minute
)

type activityStore interface {
	UpdateLastActive(ctx context.Context, userID uuid.UUID) error
}

type Tracker struct {
	client *redis.Client
	store  activityStore
}

func NewTracker(client *redis.Client, store activityStore) *Tracker {
	return &Tracker{
		client: client,
		store:  store,
	}
}

func onlineKey(userID uuid.UUID) string {
	return "presence:online:" + userID.String()
}

func syncKey(userID uuid.UUID) string {
	return "presence:sync:" + userID.String()
}

func (t *Tracker) Touch(ctx context.Context, userID uuid.UUID) {
	pipe := t.client.Pipeline()
	pipe.Set(ctx, onlineKey(userID), "1", onlineTTL)
	due := pipe.SetNX(ctx, syncKey(userID), "1", syncInterval)

	if _, err := pipe.Exec(ctx); err != nil {
		slog.WarnContext(ctx, "touch presence", "user_id", userID, "error", err)
		return
	}

	if !due.Val() {
		return
	}

	if err := t.store.UpdateLastActive(ctx, userID); err != nil {
		slog.WarnContext(ctx, "update last active", "user_id", userID, "error", err)
	}
}

func (t *Tracker) Remove(ctx context.Context, userID uuid.UUID) {
	if err := t.client.Del(ctx, onlineKey(userID)).Err(); err != nil {
		slog.WarnContext(ctx, "remove presence", "user_id", userID, "error", err)
	}
}

func (t *Tracker) Online(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	online := make(map[uuid.UUID]bool, len(userIDs))
	if len(userIDs) == 0 {
		return online, nil
	}

	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = onlineKey(id)
	}

	values, err := t.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("presence: online: %w", err)
	}

	for i, value := range values {
		online[userIDs[i]] = value != nil
	}

	return online, nil
}
