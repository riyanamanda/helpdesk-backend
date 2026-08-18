package websocket

import (
	"encoding/json"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			slog.Error(
				"websocket invalid message",
				"error", err,
			)
			continue
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		c.hub.broadcast <- payload
	}
}

func (c *Client) writePump() {
	defer func() {
		_ = c.conn.Close()
	}()

	for message := range c.send {
		if err := c.conn.WriteMessage(
			websocket.TextMessage,
			message,
		); err != nil {
			return
		}
	}
}
