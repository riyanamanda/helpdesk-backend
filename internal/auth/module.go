package auth

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	goredis "github.com/redis/go-redis/v9"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type redisAdapter struct {
	client *goredis.Client
}

func (r *redisAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisAdapter) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func Register(e *echo.Group, db *sqlx.DB, cfg config.Auth, storageConfig config.Storage, redisClient *goredis.Client) {
	repo := user.NewUserRepository(db)
	svc := NewAuthService(repo, cfg, storageConfig, &redisAdapter{client: redisClient})
	handler := NewAuthHandler(svc)

	authGroup := e.Group("/auth")
	protected := authGroup.Group("")
	protected.Use(
		middleware.AuthMiddleware(cfg, redisClient),
	)

	authGroup.POST("/login", handler.Login)
	authGroup.POST("/google", handler.LoginWithGoogle)

	protected.POST("/logout", handler.Logout)
	protected.GET("/me", handler.Me)
}
