package category

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo, cache)
	handler := NewCategoryHandler(svc)

	e.GET("/categories", handler.ListCategories)
	e.GET("/categories/options", handler.ListCategoryOptions)
	e.GET("/categories/:id", handler.GetCategory)
	e.POST("/categories", handler.CreateCategory, middleware.RequiredPermission("category:create"))
	e.PATCH("/categories/:id", handler.UpdateCategory, middleware.RequiredPermission("category:update"))
	e.DELETE("/categories/:id", handler.DeleteCategory, middleware.RequiredPermission("category:delete"))
}
