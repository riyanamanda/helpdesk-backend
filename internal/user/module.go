package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

func Register(e *echo.Group, db *sqlx.DB, storage storage.Storage, storageConfig config.Storage) {
	repo := NewUserRepository(db)
	svc := NewUserService(repo, storage, storageConfig)
	handler := NewUserHandler(svc)

	e.GET("/users", handler.ListUsers)
	e.POST("/users", handler.CreateUser)
	e.GET("/users/:id", handler.GetUser)
	e.PATCH("/users/me/avatar", handler.UpdateUserAvatar)
}
