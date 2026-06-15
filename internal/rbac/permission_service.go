package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type permissionService struct {
	repo  RBACRepository
	cache cache.Cache
}

func NewPermissionService(repo RBACRepository, cache cache.Cache) ctxkey.PermissionService {
	return &permissionService{
		repo:  repo,
		cache: cache,
	}
}

func (s *permissionService) GetUserPermissions(ctx context.Context, userID uuid.UUID) (ctxkey.PermissionSet, error) {
	cacheKey := BuildUserPermissionsCacheKey(userID)

	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var codes []string

		if err := json.Unmarshal([]byte(cached), &codes); err == nil {
			permissions := make(ctxkey.PermissionSet)
			for _, code := range codes {
				permissions[code] = struct{}{}
			}

			return permissions, nil
		}
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	permissions, err := s.repo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make(ctxkey.PermissionSet)
	codes := make([]string, 0, len(permissions))

	for _, permission := range permissions {
		result[permission.Code] = struct{}{}
		codes = append(codes, permission.Code)
	}

	if payload, err := json.Marshal(codes); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(payload), time.Hour)
	}

	return result, nil
}
