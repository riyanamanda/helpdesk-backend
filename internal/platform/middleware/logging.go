package middleware

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

func requestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			id := uuid.NewString()
			c.Set("request_id", id)
			c.Response().Header().Set("X-Request-ID", id)

			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			status := 0
			if err != nil {
				status = apperr.As(err).Status
			} else if resp, unwrapErr := echo.UnwrapResponse(c.Response()); unwrapErr == nil {
				status = resp.Status
			}

			slog.Info("request",
				"request_id", id,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", status,
				"latency_ms", latency.Milliseconds(),
			)

			return err
		}
	}
}
