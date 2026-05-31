package division

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func InvalidateCache(ctx context.Context, c cache.Cache) {
	_ = c.Delete(ctx, DivisionOptionsCacheKey)
}
