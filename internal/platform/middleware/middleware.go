package middleware

import (
	"github.com/labstack/echo/v5"
	echomw "github.com/labstack/echo/v5/middleware"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
)

func Register(e *echo.Echo, cfg config.App) {
	e.Use(echomw.Recover())
	e.Use(corsMiddleware(cfg.CORSOrigins))
	e.Use(requestID())
}
