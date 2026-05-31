package category

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Invalidate(ctx context.Context, c cache.Cache) {
	_ = c.Delete(ctx, CategoryOptionsCacheKey)
}
