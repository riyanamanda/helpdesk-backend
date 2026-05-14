package validation

import (
	"fmt"
	"mime/multipart"

	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

func ValidateImage(header *multipart.FileHeader, maxSize int64, allowedTypes map[string]bool) error {
	if header == nil {
		return apperror.BadRequest("image is required")
	}

	if header.Size > maxSize {
		message := fmt.Sprintf("image size should not exceed %d", maxSize)
		return apperror.BadRequest(message)
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return apperror.BadRequest("invalid image format")
	}

	return nil
}
