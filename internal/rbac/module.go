package rbac

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, svc RBACService) {
	handler := NewRBACHandler(svc)

	e.GET("/roles", handler.ListRoles, middleware.RequirePermission(PermissionRBACView))
	e.GET("/permissions", handler.ListPermissions, middleware.RequirePermission(PermissionRBACView))
	e.GET("/roles/:id/permissions", handler.GetRolePermissions, middleware.RequirePermission(PermissionRBACView))
	e.PUT("/roles/:id/permissions", handler.SetRolePermissions, middleware.RequirePermission(PermissionRBACUpdate))
}
