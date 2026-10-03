package bpjs

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, svc service) {
	handler := NewHandler(svc)

	e.GET("/bpjs/peserta/:nik", handler.GetPesertaByNIK, middleware.RequirePermission(rbac.PermissionIHSView))
}
