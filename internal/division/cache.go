package division

import (
	"context"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	if err := c.Delete(ctx, DivisionOptionsCacheKey); err != nil {
		slog.WarnContext(ctx, "invalidate division options cache", "error", err)
	}
}
