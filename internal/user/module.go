package user

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Register(e *echo.Group, repo UserRepository, storageConfig config.Storage, cache cache.Cache) {
	svc := NewUserService(repo, storageConfig, cache)
	handler := NewUserHandler(svc)

	e.GET("/users", handler.ListUsers)
	e.POST("/users", handler.CreateUser)
	e.GET("/users/:id", handler.GetUser)
	e.PATCH("/users/:id", handler.UpdateUser)
	e.PATCH("/users/:id/password", handler.UpdatePassword)
	e.GET("/users/assignable", handler.ListAssignableUser)
}
