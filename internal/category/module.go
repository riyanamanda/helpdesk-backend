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

	e.GET("/categories/:id", handler.GetCategory)

	e.POST("/categories", handler.CreateCategory)

	e.PATCH("/categories/:id", handler.UpdateCategory)

	e.DELETE("/categories/:id", handler.DeleteCategory)

}
