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

	adminOnly := middleware.RequireRole("ADMIN")

	e.GET("/users", handler.ListUsers, adminOnly)
	e.POST("/users", handler.CreateUser, adminOnly)
	e.GET("/users/assignable", handler.ListAssignableUser)
	e.GET("/users/:id", handler.GetUser, adminOnly)
	e.PATCH("/users/:id", handler.UpdateUser, adminOnly)
	e.PATCH("/users/:id/password", handler.UpdatePassword, adminOnly)
}
