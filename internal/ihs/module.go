package ihs

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewPatientRepository(db)
	service := NewPatientService(repo)
	handler := NewPatientHandler(service)

	e.GET("/patients", handler.ListPatients, middleware.RequirePermission(rbac.PermissionIHSView))
	e.GET("/patients/:norm/detail", handler.GetPatientDetail, middleware.RequirePermission(rbac.PermissionIHSView))
	e.PATCH("/patients/:norm", handler.UpdatePatientMethod, middleware.RequirePermission(rbac.PermissionIHSUpdate))
}
