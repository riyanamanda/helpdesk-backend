package auth

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateResetPasswordCache(ctx context.Context, c cache.Cache, userID string, token string) {
	_ = c.DeleteMany(
		ctx,
		buildPasswordResetRequestCache(userID),
		buildPasswordResetTokenCache(token),
	)
}

func buildPasswordResetCooldownCache(userID string) string {
	return "password-reset:cooldown:" + userID
}

func buildPasswordResetRequestCache(userID string) string {
	return "password-reset:request:" + userID
}

func buildPasswordResetTokenCache(token string) string {
	return "password-reset:token:" + token
}
