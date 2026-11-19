package auth

import (
	"context"
	"log/slog"
	"strings"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func InvalidateResetPasswordCache(ctx context.Context, c cache.Cache, userID string, token string) {
	err := c.DeleteMany(
		ctx,
		buildPasswordResetRequestCache(userID),
		buildPasswordResetTokenCache(token),
	)
	if err != nil {
		slog.WarnContext(ctx, "invalidate password reset token", "user_id", userID, "error", err)
	}
}

func buildLoginAttemptsCache(email string) string {
	return "auth:login:attempts:" + strings.ToLower(strings.TrimSpace(email))
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
