package ctxkey

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	ContextUserIDKey     contextKey = "user_id"
	ContextJTIKey        contextKey = "jti"
	ContextPermissionKey contextKey = "permissions"
)

func SetUserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextUserIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(ContextUserIDKey).(uuid.UUID)
	return userID, ok
}

func SetJTIToContext(ctx context.Context, jti string) context.Context {
	return context.WithValue(ctx, ContextJTIKey, jti)
}

func GetJTIFromContext(ctx context.Context) (string, bool) {
	jti, ok := ctx.Value(ContextJTIKey).(string)
	return jti, ok
}

func SetPermissionToContext(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, ContextPermissionKey, permissions)
}

func GetPermissionFromContext(ctx context.Context) ([]string, bool) {
	permissions, ok := ctx.Value(ContextPermissionKey).([]string)
	return permissions, ok
}
