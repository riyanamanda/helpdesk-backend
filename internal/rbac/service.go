package rbac

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

type repository interface {
	GetRoles(ctx context.Context) ([]Role, error)
	GetPermissions(ctx context.Context) ([]Permission, error)
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]Permission, error)
	SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
	GetUserIDsByRoleID(ctx context.Context, roleID int64) ([]uuid.UUID, error)
}

type Service struct {
	repo  repository
	cache cache.Cache
}

func NewService(repo repository, cache cache.Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) ListRoles(ctx context.Context) ([]RoleResponse, error) {
	roles, err := s.repo.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	return toRoleResponses(roles), nil
}

func (s *Service) ListPermissions(ctx context.Context) ([]PermissionResponse, error) {
	permissions, err := s.repo.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return toPermissionResponses(permissions), nil
}

func (s *Service) GetRolePermissions(ctx context.Context, roleID int64) ([]PermissionResponse, error) {
	permissions, err := s.repo.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return toPermissionResponses(permissions), nil
}

func (s *Service) SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
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
