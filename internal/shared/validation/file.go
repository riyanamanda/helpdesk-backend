package validation

import (
	"mime/multipart"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

func ValidateImage(header *multipart.FileHeader, maxSize int64, allowedTypes map[string]bool) error {
	if header == nil {
		return nil
	}

	if header.Size > maxSize {
		return apperr.BadRequest("file is too large")
	}

	contentType := header.Header.Get("Content-Type")

	if !allowedTypes[contentType] {
		return apperr.BadRequest("invalid image format")
	}

	return nil
}
