package middleware

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
)

func AuthMiddleware(cfg config.Auth) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, apperror.Unauthorized(apperror.CodeMissingToken, "missing authorization header"))
			}

			const bearerPrefix = "Bearer "

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return response.Error(c, apperror.Unauthorized(apperror.CodeMissingToken, "invalid authorization header"))
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := utils.ParseToken(tokenString, cfg.JWTSecret)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					return response.Error(c, apperror.Unauthorized(apperror.CodeTokenExpired, "token expired"))
				}

				return response.Error(c, apperror.Unauthorized(apperror.CodeInvalidToken, "invalid token"))
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return response.Error(c, apperror.Unauthorized(apperror.CodeInvalidToken, "invalid token"))
			}

			ctx := utils.SetUserIDToContext(c.Request().Context(), userID)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
