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

	roleCode, err := s.getUserRoleCode(ctx, userID)
	if err != nil {
		return nil, err
	}

	var rawPerms []Permission
	if roleCode == string(SUPERADMIN) {
		rawPerms, err = s.repo.GetPermissions(ctx)
	} else {
		rawPerms, err = s.repo.GetPermissionsByUserID(ctx, userID)
	}
	if err != nil {
		return nil, err
	}

	result := make(ctxkey.PermissionSet)
	codes := make([]string, 0, len(rawPerms))
	for _, p := range rawPerms {
		result[p.Code] = struct{}{}
		codes = append(codes, p.Code)
	}

	if payload, err := json.Marshal(codes); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(payload), time.Hour)
	}

	return result, nil
}

func (s *permissionService) getUserRoleCode(ctx context.Context, userID uuid.UUID) (string, error) {
	cacheKey := BuildUserRoleCacheKey(userID)

	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		return cached, nil
	}

	code, err := s.repo.GetUserRoleCode(ctx, userID)
	if err != nil {
		return "", err
	}

	_ = s.cache.Set(ctx, cacheKey, code, time.Hour)
	return code, nil
}
