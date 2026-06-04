package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := ctxkey.GetRoleFromContext(c.Request().Context())
			if !ok || !allowed[role] {
				return response.Error(c, apperr.Forbidden("access denied"))
			}
			return next(c)
		}
	}
}
