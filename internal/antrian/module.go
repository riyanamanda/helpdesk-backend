package antrian

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, svc AntrianService) {
	handler := NewAntrianHandler(svc)

	e.GET("/antrian", handler.ListAntrian, middleware.RequirePermission(rbac.PermissionAntrianView))
	e.POST("/antrian/:kode_booking/checkin", handler.CheckIn, middleware.RequirePermission(rbac.PermissionAntrianCheckIn))
}
