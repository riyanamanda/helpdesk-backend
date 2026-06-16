package rbac

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func BuildUserPermissionsCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf(UserPermissionsCacheKey, userID.String())
}

func InvalidateUserPermissionsCache(ctx context.Context, c cache.Cache, userID uuid.UUID) {
	_ = c.Delete(ctx, BuildUserPermissionsCacheKey(userID))
}

func BuildUserRoleCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf(UserRoleCacheKey, userID.String())
}

func InvalidateUserRoleCache(ctx context.Context, c cache.Cache, userID uuid.UUID) {
	_ = c.Delete(ctx, BuildUserRoleCacheKey(userID))
}
