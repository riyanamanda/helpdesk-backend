package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

func RequirePermission(permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authUser, ok := ctxkey.GetAuthUserFromContext(c.Request().Context())
			if !ok || authUser == nil {
				return apperr.Forbidden("forbidden")
			}

			if !authUser.Permissions.Has(permissions...) {
				return apperr.Forbidden("forbidden")
			}

			return next(c)
		}
	}
}
