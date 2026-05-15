package validation

import (
	"mime/multipart"

	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

func ValidateImage(header *multipart.FileHeader, maxSize int64, allowedTypes map[string]bool) error {
	if header == nil {
		return nil
	}

	if header.Size > maxSize {
		return apperror.BadRequest("file is too large")
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return apperror.BadRequest("invalid image format")
	}

	return nil
}
