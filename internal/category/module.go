package category

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo, cache)
	handler := NewCategoryHandler(svc)

	e.GET("/categories", handler.ListCategories)
	e.GET("/categories/options", handler.ListCategoryOption)
	e.GET("/categories/:id", handler.GetCategory)
	e.POST("/categories", handler.CreateCategory)
	e.PATCH("/categories/:id", handler.UpdateCategory)
	e.DELETE("/categories/:id", handler.DeleteCategory)
}
