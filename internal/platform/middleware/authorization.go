package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

func RequiredPermission(permissions ...string) echo.MiddlewareFunc {
	allowed := make(map[string]bool, len(permissions))
	for _, p := range permissions {
		allowed[p] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userPermissions, ok := ctxkey.GetPermissionFromContext(c.Request().Context())
			if !ok {
				return response.Error(c, apperr.Forbidden("access denied"))
			}

			for _, p := range userPermissions {
				if allowed[p] {
					return next(c)
				}
			}

			return response.Error(c, apperr.Forbidden("access denied"))
		}
	}
}
