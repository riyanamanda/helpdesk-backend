package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(e *echo.Group, db *sqlx.DB, cfg config.Auth, storageConfig config.Storage) {
	repo := user.NewUserRepository(db)
	svc := NewAuthService(repo, cfg, storageConfig)
	handler := NewAuthHandler(svc)

	authGroup := e.Group("/auth")
	protected := authGroup.Group("")
	protected.Use(
		middleware.AuthMiddleware(cfg),
	)

	authGroup.POST("/login", handler.Login)
	authGroup.POST("/google", handler.LoginWithGoogle)

	protected.GET("/me", handler.Me)
}
