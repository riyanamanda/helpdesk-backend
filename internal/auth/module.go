package auth

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(e *echo.Group, db *sqlx.DB, jwtSecret string, jwtExpirationMinutes int) {
	repo := user.NewUserRepository(db)
	svc := NewAuthService(repo, jwtSecret, time.Duration(jwtExpirationMinutes)*time.Minute)
	handler := NewAuthHandler(svc)

	e.POST("/auth/login", handler.Login)
}
