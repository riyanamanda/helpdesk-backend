package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
)

func AssertAppError(t *testing.T, err error, code string, status int, message string) {

	t.Helper()

	var appErr *apperror.AppError

	require.ErrorAs(t, err, &appErr)

	assert.Equal(t, code, appErr.Code)

	assert.Equal(t, status, appErr.StatusCode)

	assert.Equal(t, message, appErr.Message)

}
