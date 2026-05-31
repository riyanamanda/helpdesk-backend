package dashboard

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	_ = c.Delete(ctx, SummaryCacheKey)
	_ = c.Delete(ctx, RecentTicketsCacheKey)
}
