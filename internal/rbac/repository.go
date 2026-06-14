package rbac

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RBACRepository interface {
	GetAllRoles(ctx context.Context) ([]Role, error)

	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetAllPermissionsByRoleID(ctx context.Context, roleID int64) ([]Permission, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRBACRepository(db *sqlx.DB) RBACRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAllRoles(ctx context.Context) ([]Role, error) {
	var roles []Role

	const query = `
		SELECT id, code
		FROM roles
	`

	if err := r.db.SelectContext(ctx, &roles, query); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *repository) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	var permissions []Permission

	const query = `
		SELECT id, code
		FROM permissions
	`

	if err := r.db.SelectContext(ctx, &permissions, query); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *repository) GetAllPermissionsByRoleID(ctx context.Context, roleID int64) ([]Permission, error) {
	var permissions []Permission

	const query = `
		SELECT
			p.id as id,
			p.code as code
		FROM role_permissions rp
		JOIN roles r
			ON r.id = rp.role_id
		JOIN permissions p
			ON p.id = rp.premission_id
		WHERE role_id = $1
	`

	if err := r.db.SelectContext(ctx, &permissions, query, roleID); err != nil {
		return nil, err
	}

	return permissions, nil
}
