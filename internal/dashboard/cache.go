package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	_ = c.Delete(ctx, SummaryCacheKey)
	_ = c.Delete(ctx, RecentTicketsCacheKey)
	_ = c.Delete(ctx, CategoryTicketsCacheKey)
	_ = c.Delete(ctx, AgentWorkloadCacheKey)

	currentYear := time.Now().Year()
	for _, year := range []int{currentYear - 1, currentYear} {
		_ = c.Delete(ctx, fmt.Sprintf(MonthlyTrendCacheKey, year))
	}
}
