package rbac

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

type RBACService interface {
	ListRoles(ctx context.Context) ([]RoleResponse, error)
	ListPermissions(ctx context.Context) ([]PermissionResponse, error)
	GetRolePermissions(ctx context.Context, roleID int64) ([]PermissionResponse, error)
	SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
}

type service struct {
	repo  RBACRepository
	cache cache.Cache
}

func NewRBACService(repo RBACRepository, cache cache.Cache) RBACService {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) ListRoles(ctx context.Context) ([]RoleResponse, error) {
	roles, err := s.repo.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	return toRoleResponses(roles), nil
}

func (s *service) ListPermissions(ctx context.Context) ([]PermissionResponse, error) {
	permissions, err := s.repo.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return toPermissionResponses(permissions), nil
}

func (s *service) GetRolePermissions(ctx context.Context, roleID int64) ([]PermissionResponse, error) {
	permissions, err := s.repo.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return toPermissionResponses(permissions), nil
}

func (s *service) SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if err := s.repo.SetRolePermissions(ctx, roleID, permissionIDs); err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			return apperr.BadRequest("one or more permission IDs are invalid")
		}
		return err
	}

	userIDs, err := s.repo.GetUserIDsByRoleID(ctx, roleID)
	if err != nil {
		return err
	}

	keys := make([]string, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = BuildUserPermissionsCacheKey(userID)
	}
	_ = s.cache.DeleteMany(ctx, keys...)

	return nil
}
