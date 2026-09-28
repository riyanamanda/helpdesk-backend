package ihs

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, svc PatientService) {
	handler := NewPatientHandler(svc)

	e.GET("/patients", handler.ListPatients, middleware.RequirePermission(rbac.PermissionIHSView))
	e.GET("/patients/:norm/detail", handler.GetPatientDetail, middleware.RequirePermission(rbac.PermissionIHSView))
	e.PATCH("/patients/:norm", handler.UpdatePatientMethod, middleware.RequirePermission(rbac.PermissionIHSUpdate))
}
