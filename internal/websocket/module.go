package websocket

import "github.com/labstack/echo/v5"

func Register(e *echo.Group, hub *Hub) {
	handler := NewHandler(hub)

	e.GET("/ws", handler.Connect)
}
