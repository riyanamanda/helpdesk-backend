package category

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo)
	handler := NewCategoryHandler(svc)

	e.GET("/categories", handler.ListCategories)
	e.GET("/categories/:id", handler.GetByID)
	e.POST("/categories", handler.Create)
	e.PUT("/categories/:id", handler.Update)
	e.DELETE("/categories/:id", handler.Delete)
}
