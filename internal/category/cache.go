package category

import (
	"context"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	if err := c.Delete(ctx, CategoryOptionsCacheKey); err != nil {
		slog.WarnContext(ctx, "invalidate category options cache", "error", err)
	}
}
