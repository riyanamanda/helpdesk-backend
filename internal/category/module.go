package category

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo, cache)
	handler := NewCategoryHandler(svc)

	adminOnly := middleware.RequireRole("ADMIN")

	e.GET("/categories", handler.ListCategories)
	e.GET("/categories/options", handler.ListCategoryOptions)
	e.GET("/categories/:id", handler.GetCategory)
	e.POST("/categories", handler.CreateCategory, adminOnly)
	e.PATCH("/categories/:id", handler.UpdateCategory, adminOnly)
	e.DELETE("/categories/:id", handler.DeleteCategory, adminOnly)
}
