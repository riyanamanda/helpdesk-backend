package middleware

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/jwtutil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

func AuthMiddleware(cfg config.Auth, redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, apperr.Unauthorized(apperr.CodeMissingToken, "missing authorization header"))
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return response.Error(c, apperr.Unauthorized(apperr.CodeMissingToken, "invalid authorization header"))
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := jwtutil.ParseToken(tokenString, cfg.JWTSecret)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					return response.Error(c, apperr.Unauthorized(apperr.CodeTokenExpired, "token expired"))
				}

				return response.Error(c, apperr.Unauthorized(apperr.CodeInvalidToken, "invalid token"))
			}

			key := jwtutil.TokenKeyPrefix + claims.ID
			exists, err := redisClient.Exists(c.Request().Context(), key).Result()
			if err != nil {
				return response.Error(c, apperr.Internal())
			}

			if exists == 0 {
				return response.Error(c, apperr.Unauthorized(apperr.CodeInvalidToken, "token revoked"))
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return response.Error(c, apperr.Unauthorized(apperr.CodeInvalidToken, "invalid token"))
			}

			ctx := ctxkey.SetUserIDToContext(c.Request().Context(), userID)
			ctx = ctxkey.SetJTIToContext(ctx, claims.ID)
			ctx = ctxkey.SetRoleToContext(ctx, claims.Role)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
