package main

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/antrian"
	"github.com/riyanamanda/helpdesk-backend/internal/auth"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/ihs"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/profile"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
	"github.com/riyanamanda/helpdesk-backend/internal/ticket"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func registerRoutes(cfg *config.Config, d *deps) http.Handler {
	e := echo.New()
	e.Validator = validation.New()
	middleware.Register(e, cfg.App)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"name":   cfg.App.Name,
		})
	})

	e.Match([]string{http.MethodGet, http.MethodHead}, "/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")
	auth.Register(api, d.userRepo, cfg.Auth, cfg.Storage, d.redisClient, d.permissionService)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.Auth, d.redisClient, d.permissionService))

	category.Register(protected, d.db, d.cacheStore)
	division.Register(protected, d.db, d.cacheStore)
	user.Register(protected, d.userRepo, d.txManager, cfg.Storage, d.cacheStore)
	ticket.Register(protected, d.db, d.storageService, cfg.Storage, d.cacheStore, d.userRepo)
	dashboard.Register(protected, d.db, d.cacheStore)
	profile.Register(protected, d.db, d.storageService, cfg.Storage, cfg.Auth)
	rbac.Register(protected, d.db, d.cacheStore)
	if d.simgosDB != nil {
		ihs.Register(protected, d.simgosDB)
		antrian.Register(protected, d.simgosDB, cfg.Antrol)
	}

	return e
}
