package middleware

import (
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
)

func Register(e *echo.Echo, cfg config.App) {

	e.Use(corsMiddleware(cfg.CORSOrigins))

	e.Use(requestID())

}
