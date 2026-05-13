package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"

	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
)

func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, apperror.Forbidden("missing authorization header"))
			}

			const bearerPrefix = "Bearer "

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return response.Error(c, apperror.Forbidden("invalid authorization header"))
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := utils.ParseToken(tokenString, jwtSecret)
			if err != nil {
				return response.Error(c, apperror.Forbidden("invalid token"))
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return response.Error(c, apperror.Forbidden("invalid token"))
			}

			ctx := utils.SetUserIDToContext(c.Request().Context(), userID)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
