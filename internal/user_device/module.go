package user_device

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewUserDeviceRepository(db)
	svc := NewUserDeviceService(repo)
	handler := NewUserDeviceHandler(svc)

	deviceGroup := e.Group("/devices")
	deviceGroup.POST("", handler.RegisterDevice)
	deviceGroup.DELETE("", handler.UnregisterDevice)
}
