package ihs

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, db *sqlx.DB, cfg config.Database) {
	repo := NewPatientRepository(db)
	simgos := newSimgosClient(cfg.Host)
	service := NewPatientService(repo, simgos)
	handler := NewPatientHandler(service)

	e.GET("/patients", handler.ListPatients, middleware.RequirePermission(rbac.PermissionIHSView))
	e.GET("/patients/:norm/detail", handler.GetPatientDetail, middleware.RequirePermission(rbac.PermissionIHSView))
	e.PATCH("/patients/:norm", handler.UpdatePatientMethod, middleware.RequirePermission(rbac.PermissionIHSUpdate))
	e.GET("/ihs/patient/send", handler.SendIhs, middleware.RequirePermission(rbac.PermissionIHSUpdate))
}
