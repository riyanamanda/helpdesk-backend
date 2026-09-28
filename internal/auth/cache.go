package auth

import (
	"context"
	"fmt"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateResetPasswordCache(ctx context.Context, c cache.Cache, userID string, token string) {
	_ = c.Delete(ctx, buildPasswordResetRequestCache(userID))
	_ = c.Delete(ctx, buildPasswordResetCache(token))
}

func buildPasswordResetRequestCache(userID string) string {
	return fmt.Sprintf(PasswordResetRequestCacheKey, userID)
}

func buildPasswordResetCache(token string) string {
	return fmt.Sprintf(PasswordResetCacheKey, token)
}
