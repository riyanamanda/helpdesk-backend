package passwordtoken

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

func New() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate password token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func CacheKey(token string) string {
	return "password-reset:token:" + token
}

func URL(appURL, token string) string {
	return fmt.Sprintf("%s/reset-password?token=%s", appURL, url.QueryEscape(token))
}

func Issue(ctx context.Context, c cache.Cache, appURL string, userID uuid.UUID, ttl time.Duration) (string, error) {
	token, err := New()
	if err != nil {
		return "", err
	}

	if err := c.Set(ctx, CacheKey(token), userID.String(), ttl); err != nil {
		return "", fmt.Errorf("store password token: %w", err)
	}

	return URL(appURL, token), nil
}
