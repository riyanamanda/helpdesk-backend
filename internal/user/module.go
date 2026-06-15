package user

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, repo UserRepository, storageConfig config.Storage, cache cache.Cache) {
	svc := NewUserService(repo, storageConfig, cache)
	handler := NewUserHandler(svc)

	e.GET("/users", handler.ListUsers, middleware.RequiredPermission("user:view"))
	e.POST("/users", handler.CreateUser, middleware.RequiredPermission("user:create"))
	e.GET("/users/assignable", handler.ListAssignableUser)
	e.GET("/users/:id", handler.GetUser, middleware.RequiredPermission("user:view"))
	e.PATCH("/users/:id", handler.UpdateUser, middleware.RequiredPermission("user:update"))
	e.PATCH("/users/:id/password", handler.UpdatePassword, middleware.RequiredPermission("user:update"))
}
