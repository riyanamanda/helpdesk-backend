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
	e.POST("/categories", handler.Create)
}
