package middleware

import "github.com/labstack/echo/v5"

func Register(e *echo.Echo) {
	registerCore(e)
	registerObservability(e)
}

func registerCore(e *echo.Echo) {
	e.Use(recoverMiddleware())
}

func registerObservability(e *echo.Echo) {
	e.Use(requestID())
	// e.Use(requestLogger())
}
