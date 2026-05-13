package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewUserRepository(db)
	svc := NewUserService(repo)
	handler := NewUserHandler(svc)

	e.GET("/users", handler.ListUser)
	e.POST("/users", handler.Create)
	e.GET("/users/:id", handler.GetByID)
}
