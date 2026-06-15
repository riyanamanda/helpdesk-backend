package category

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo, cache)
	handler := NewCategoryHandler(svc)

	e.GET("/categories", handler.ListCategories, middleware.RequirePermission(rbac.PermissionCategoryView))
	e.GET("/categories/options", handler.ListCategoryOptions)
	e.GET("/categories/:id", handler.GetCategory, middleware.RequirePermission(rbac.PermissionCategoryView))
	e.POST("/categories", handler.CreateCategory, middleware.RequirePermission(rbac.PermissionCategoryCreate))
	e.PATCH("/categories/:id", handler.UpdateCategory, middleware.RequirePermission(rbac.PermissionCategoryUpdate))
	e.DELETE("/categories/:id", handler.DeleteCategory, middleware.RequirePermission(rbac.PermissionCategoryDelete))
}
