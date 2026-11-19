package user

import (
	"context"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	if err := c.Delete(ctx, AssignableCacheKey); err != nil {
		slog.WarnContext(ctx, "invalidate assignable users cache", "error", err)
	}
}
