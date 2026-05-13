package utils

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func GetUserIDFromJWT(c *echo.Context) (uuid.UUID, bool) {
	claims, ok := c.Get("user").(*JWTCustomClaims)
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}
