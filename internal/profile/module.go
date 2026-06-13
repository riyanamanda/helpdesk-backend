package profile

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
)

func Register(e *echo.Group, db *sqlx.DB, storageService storage.Storage, storageConfig config.Storage, authConfig config.Auth) {
	repo := NewProfileRepository(db)
	svc := NewProfileService(repo, storageService, storageConfig, authConfig)
	handler := NewProfileHandler(svc)

	profileGroup := e.Group("/profile")

	profileGroup.GET("", handler.GetProfile)
	profileGroup.PUT("", handler.UpdateProfile)
	profileGroup.PATCH("/avatar", handler.UpdateAvatar)
	profileGroup.POST("/sync-google", handler.SyncGoogle)
	profileGroup.POST("/revoke-google", handler.RevokeGoogle)
	profileGroup.PATCH("/update-password", handler.UpdatePassword)
}
