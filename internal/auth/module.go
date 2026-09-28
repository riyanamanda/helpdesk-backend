package auth

import (
	"github.com/labstack/echo/v5"
	goredis "github.com/redis/go-redis/v9"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(
	e *echo.Group, userRepo user.UserRepository, cfg config.Auth, storageConfig config.Storage,
	redisClient *goredis.Client, permissionService ctxkey.PermissionService, rabbitmq rabbitmq.Client,
) {
	svc := NewAuthService(userRepo, cfg, storageConfig, cache.NewRedisCache(redisClient), permissionService, rabbitmq)
	handler := NewAuthHandler(svc)

	authGroup := e.Group("/auth")
	protected := authGroup.Group("")

	protected.Use(
		middleware.AuthMiddleware(cfg, redisClient, permissionService),
	)

	authGroup.POST("/login", handler.Login)
	authGroup.POST("/google", handler.LoginWithGoogle)
	authGroup.POST("/forgot-password", handler.ForgotPassword)

	protected.POST("/logout", handler.Logout)
	protected.GET("/me", handler.Me)
}
