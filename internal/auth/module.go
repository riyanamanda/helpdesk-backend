package auth

import (
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, svc AuthService, authMiddleware echo.MiddlewareFunc) {
	handler := NewAuthHandler(svc)

	e.POST("/auth/login", handler.Login)
	e.POST("/auth/google", handler.LoginWithGoogle)
	e.POST("/auth/forgot-password", handler.ForgotPassword)
	e.POST("/auth/reset-password", handler.ResetPassword)

	protected := e.Group("")
	protected.Use(authMiddleware)

	protected.Use()
	protected.POST("/auth/logout", handler.Logout)
	protected.GET("/auth/me", handler.Me)
}
