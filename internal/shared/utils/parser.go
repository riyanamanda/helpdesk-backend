package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

func ParsePositiveInt64PathParam(c *echo.Context, paramName, resourceName string) (int64, error) {
	value := c.Param(paramName)
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		if resourceName != "" {
			return 0, apperrors.BadRequest(fmt.Sprintf("invalid %s %s", resourceName, paramName))
		}
		return 0, apperrors.BadRequest(fmt.Sprintf("invalid %s", paramName))
	}

	return id, nil
}

func GenerateObjectKey(prefix, filename string) string {
	ext := filepath.Ext(filename)

	random := strings.ReplaceAll(
		uuid.NewString(),
		"-",
		"",
	)[:8]

	return fmt.Sprintf(
		"%s/%d-%s%s",
		prefix,
		time.Now().Unix(),
		random,
		ext,
	)
}
