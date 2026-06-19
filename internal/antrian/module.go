package antrian

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, db *sqlx.DB, cfg config.Antrol) {
	repo := NewAntrianRepository(db)
	antrol := newAntrolClient(cfg.Domain, cfg.Username, cfg.Password)
	svc := NewAntrianService(repo, antrol)
	h := NewAntrianHandler(svc)

	e.GET("/antrian", h.ListAntrian, middleware.RequirePermission(rbac.PermissionAntrianView))
	e.POST("/antrian/:kode_booking/checkin", h.CheckIn, middleware.RequirePermission(rbac.PermissionAntrianCheckIn))
}
