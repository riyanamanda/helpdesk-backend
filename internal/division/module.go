package division

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewDivisionRepository(db)
	svc := NewDivisionService(repo, cache)
	handler := NewDivisionHandler(svc)

	e.GET("/divisions", handler.ListDivisions)
	e.GET("/divisions/options", handler.ListDivisionOptions)
	e.GET("/divisions/:id", handler.GetDivision)
	e.POST("/divisions", handler.CreateDivision)
	e.PATCH("/divisions/:id", handler.UpdateDivision)
	e.DELETE("/divisions/:id", handler.DeleteDivision)
}
