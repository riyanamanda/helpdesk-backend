package rbac

import "context"

type RBACService interface {
	ListRoles(ctx context.Context) ([]RoleResponse, error)
	ListPermissions(ctx context.Context) ([]PermissionResponse, error)
}

type service struct {
	repo RBACRepository
}

func NewRBACService(repo RBACRepository) RBACService {
	return &service{
		repo: repo,
	}
}

func (s *service) ListRoles(ctx context.Context) ([]RoleResponse, error) {
	roles, err := s.repo.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	return toRoleResponses(roles), nil
}

func (s *service) ListPermissions(ctx context.Context) ([]PermissionResponse, error) {
	permissions, err := s.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return toPermissionResponses(permissions), nil
}
