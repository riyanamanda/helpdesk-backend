package ctxkey

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	contextUserIDKey   contextKey = "user_id"
	contextJTIKey      contextKey = "jti"
	contextAuthUserKey contextKey = "auth_user"
)

type PermissionSet map[string]struct{}

func (p PermissionSet) Has(permissions ...string) bool {
	for _, permission := range permissions {
		if _, ok := p[permission]; ok {
			return true
		}
	}
	return false
}

func (p PermissionSet) ToSlice() []string {
	codes := make([]string, 0, len(p))
	for code := range p {
		codes = append(codes, code)
	}
	return codes
}

type PermissionService interface {
	GetUserPermissions(ctx context.Context, userID uuid.UUID) (PermissionSet, error)
}

type AuthUser struct {
	ID          uuid.UUID
	Permissions PermissionSet
}

func SetUserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextUserIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(contextUserIDKey).(uuid.UUID)
	return userID, ok
}

func SetJTIToContext(ctx context.Context, jti string) context.Context {
	return context.WithValue(ctx, contextJTIKey, jti)
}

func GetJTIFromContext(ctx context.Context) (string, bool) {
	jti, ok := ctx.Value(contextJTIKey).(string)
	return jti, ok
}

func SetAuthUserToContext(ctx context.Context, authUser *AuthUser) context.Context {
	return context.WithValue(ctx, contextAuthUserKey, authUser)
}

func GetAuthUserFromContext(ctx context.Context) (*AuthUser, bool) {
	authUser, ok := ctx.Value(contextAuthUserKey).(*AuthUser)
	return authUser, ok
}
