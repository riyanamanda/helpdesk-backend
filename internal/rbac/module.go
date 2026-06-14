package rbac

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewRBACRepository(db)
	service := NewRBACService(repo)
	handler := NewRBACHandler(service)

	e.GET("/roles", handler.ListRoles)
	e.GET("/permissions", handler.ListPermissions)
}
