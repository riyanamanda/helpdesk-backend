package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(e *echo.Group, db *sqlx.DB, cfg config.Auth) {
	repo := user.NewUserRepository(db)
	svc := NewAuthService(repo, cfg)
	handler := NewAuthHandler(svc)

	e.POST("/auth/login", handler.Login)
}
