package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	currentYear := time.Now().Year()

	keys := []string{
		SummaryCacheKey,
		RecentTicketsCacheKey,
		CategoryTicketsCacheKey,
		AgentWorkloadCacheKey,
		fmt.Sprintf(MonthlyTrendCacheKey, currentYear-1),
		fmt.Sprintf(MonthlyTrendCacheKey, currentYear),
	}

	if err := c.DeleteMany(ctx, keys...); err != nil {
		slog.WarnContext(ctx, "invalidate dashboard cache", "error", err)
	}
}
