package rbac

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func BuildUserPermissionsCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf(UserPermissionsCacheKey, userID.String())
}

func InvalidateUserPermissionsCache(ctx context.Context, c cache.Cache, userID uuid.UUID) {
	if err := c.Delete(ctx, BuildUserPermissionsCacheKey(userID)); err != nil {
		slog.WarnContext(ctx, "invalidate user permissions cache", "user_id", userID, "error", err)
	}
}

func BuildUserRoleCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf(UserRoleCacheKey, userID.String())
}

func InvalidateUserRoleCache(ctx context.Context, c cache.Cache, userID uuid.UUID) {
	if err := c.Delete(ctx, BuildUserRoleCacheKey(userID)); err != nil {
		slog.WarnContext(ctx, "invalidate user role cache", "user_id", userID, "error", err)
	}
}
