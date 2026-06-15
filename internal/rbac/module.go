package rbac

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewRBACRepository(db)
	svc := NewRBACService(repo, cache)
	handler := NewRBACHandler(svc)

	e.GET("/roles", handler.ListRoles, middleware.RequirePermission(PermissionRBACView))
	e.GET("/permissions", handler.ListPermissions, middleware.RequirePermission(PermissionRBACView))
	e.GET("/roles/:id/permissions", handler.GetRolePermissions, middleware.RequirePermission(PermissionRBACView))
	e.PUT("/roles/:id/permissions", handler.SetRolePermissions, middleware.RequirePermission(PermissionRBACUpdate))
}
