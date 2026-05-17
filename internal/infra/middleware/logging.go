package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func requestID() echo.MiddlewareFunc {
	return middleware.RequestID()
}

func recoverMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}

// func requestLogger() echo.MiddlewareFunc {
// 	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
// 		LogRequestID: true,
// 		LogMethod:    true,
// 		LogURI:       true,
// 		LogStatus:    true,
// 		LogLatency:   true,
// 		HandleError:  true,
// 		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
// 			level := slog.LevelInfo

// 			switch {
// 			case v.Status >= 500:
// 				level = slog.LevelError
// 			case v.Status >= 400:
// 				level = slog.LevelWarn
// 			}

// 			// Log all 5xx status codes or if there's an error
// 			if v.Status >= 500 || v.Error != nil {
// 				fields := []any{
// 					"request_id", v.RequestID,
// 					"method", v.Method,
// 					"uri", v.URI,
// 					"status", v.Status,
// 					"latency", v.Latency,
// 				}
// 				if v.Error != nil {
// 					fields = append(fields, "error", v.Error)
// 				}
// 				slog.Log(c.Request().Context(), level, "http request failed", fields...)
// 				return nil
// 			}

// 			slog.Log(c.Request().Context(), level, "http request",
// 				"request_id", v.RequestID,
// 				"method", v.Method,
// 				"uri", v.URI,
// 				"status", v.Status,
// 				"latency", v.Latency,
// 			)

// 			return nil
// 		},
// 	})
// }
