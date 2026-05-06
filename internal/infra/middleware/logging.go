package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func requestID() echo.MiddlewareFunc {
	return middleware.RequestID()
}

func recoverMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}

func requestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		HandleError:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			level := slog.LevelInfo

			switch {
			case v.Status >= 500:
				level = slog.LevelError
			case v.Status >= 400:
				level = slog.LevelWarn
			}

			if v.Error != nil {
				slog.Log(c.Request().Context(), level, "http request failed",
					"request_id", v.RequestID,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"error", v.Error,
				)
				return nil
			}

			slog.Log(c.Request().Context(), level, "http request",
				"request_id", v.RequestID,
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
			)

			return nil
		},
	})
}
