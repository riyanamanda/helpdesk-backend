package division

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {

	repo := NewDivisionRepository(db)

	svc := NewDivisionService(repo)

	handler := NewDivisionHandler(svc)

	e.GET("/divisions", handler.ListDivisions)

	e.GET("/divisions/:id", handler.GetDivision)

	e.POST("/divisions", handler.CreateDivision)

	e.PATCH("/divisions/:id", handler.UpdateDivision)

	e.DELETE("/divisions/:id", handler.DeleteDivision)

}
