package user

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, repo UserRepository, storageConfig config.Storage, cache cache.Cache) {
	svc := NewUserService(repo, storageConfig, cache)
	handler := NewUserHandler(svc)

	e.GET("/users", handler.ListUsers, middleware.RequirePermission(rbac.PermissionUserView))
	e.POST("/users", handler.CreateUser, middleware.RequirePermission(rbac.PermissionUserCreate))
	e.GET("/users/assignable", handler.ListAssignableUser)
	e.GET("/users/:id", handler.GetUser, middleware.RequirePermission(rbac.PermissionUserView))
	e.PATCH("/users/:id", handler.UpdateUser, middleware.RequirePermission(rbac.PermissionUserUpdate))
	e.PATCH("/users/:id/password", handler.UpdatePassword, middleware.RequirePermission(rbac.PermissionUserUpdate))
}
