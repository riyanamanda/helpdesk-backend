package jwtutil

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func BuildUserSessionsKey(userID uuid.UUID) string {
	return fmt.Sprintf("auth:sessions:%s", userID.String())
}

func TrackSession(ctx context.Context, c cache.Cache, userID uuid.UUID, jti string, ttl time.Duration) error {
	key := BuildUserSessionsKey(userID)

	if err := c.SAdd(ctx, key, jti); err != nil {
		return err
	}

	return c.Expire(ctx, key, ttl)
}

func UntrackSession(ctx context.Context, c cache.Cache, userID uuid.UUID, jti string) error {
	return c.SRem(ctx, BuildUserSessionsKey(userID), jti)
}

func RevokeUserSessions(ctx context.Context, c cache.Cache, userID uuid.UUID) error {
	key := BuildUserSessionsKey(userID)

	jtis, err := c.SMembers(ctx, key)
	if err != nil {
		return err
	}

	if len(jtis) > 0 {
		keys := make([]string, len(jtis))
		for i, jti := range jtis {
			keys[i] = TokenKeyPrefix + jti
		}

		if err := c.DeleteMany(ctx, keys...); err != nil {
			return err
		}
	}

	return c.Delete(ctx, key)
}
