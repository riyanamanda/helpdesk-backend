package profile

import (
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, svc service) {
	handler := NewHandler(svc)

	profileGroup := e.Group("/profile")

	profileGroup.GET("", handler.GetProfile)
	profileGroup.PUT("", handler.UpdateProfile)
	profileGroup.PATCH("/avatar", handler.UpdateAvatar)
	profileGroup.POST("/sync-google", handler.SyncGoogle)
	profileGroup.POST("/revoke-google", handler.RevokeGoogle)
	profileGroup.PATCH("/update-password", handler.UpdatePassword)
}
