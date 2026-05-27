package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
)

func Register(e *echo.Group, db *sqlx.DB, storageConfig config.Storage) {
	repo := NewUserRepository(db)
	svc := NewUserService(repo, storageConfig)
	handler := NewUserHandler(svc)

	e.GET("/users", handler.ListUsers)
	e.POST("/users", handler.CreateUser)
	e.GET("/users/:id", handler.GetUser)
}
