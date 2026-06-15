package division

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewDivisionRepository(db)
	svc := NewDivisionService(repo, cache)
	handler := NewDivisionHandler(svc)

	e.GET("/divisions", handler.ListDivisions, middleware.RequirePermission(rbac.PermissionDivisionView))
	e.GET("/divisions/options", handler.ListDivisionOptions)
	e.GET("/divisions/:id", handler.GetDivision, middleware.RequirePermission(rbac.PermissionDivisionView))
	e.POST("/divisions", handler.CreateDivision, middleware.RequirePermission(rbac.PermissionDivisionCreate))
	e.PATCH("/divisions/:id", handler.UpdateDivision, middleware.RequirePermission(rbac.PermissionDivisionUpdate))
	e.DELETE("/divisions/:id", handler.DeleteDivision, middleware.RequirePermission(rbac.PermissionDivisionDelete))
}
