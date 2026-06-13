package user

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	_ = c.Delete(ctx, AssignableCacheKey)
}
