package division

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewDivisionRepository(db)
	svc := NewDivisionService(repo, cache)
	handler := NewDivisionHandler(svc)

	e.GET("/divisions", handler.ListDivisions)
	e.GET("/divisions/options", handler.ListDivisionOptions)
	e.GET("/divisions/:id", handler.GetDivision)
	e.POST("/divisions", handler.CreateDivision, middleware.RequiredPermission("division:create"))
	e.PATCH("/divisions/:id", handler.UpdateDivision, middleware.RequiredPermission("division:update"))
	e.DELETE("/divisions/:id", handler.DeleteDivision, middleware.RequiredPermission("division:delete"))
}
