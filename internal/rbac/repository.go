package rbac

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
)

type RBACRepository interface {
	GetRoles(ctx context.Context) ([]Role, error)

	GetPermissions(ctx context.Context) ([]Permission, error)
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]Permission, error)
	GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]Permission, error)

	GetUserRoleCode(ctx context.Context, userID uuid.UUID) (string, error)
	SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
	GetUserIDsByRoleID(ctx context.Context, roleID int64) ([]uuid.UUID, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRBACRepository(db *sqlx.DB) RBACRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetRoles(ctx context.Context) ([]Role, error) {
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

func (r *repository) GetPermissions(ctx context.Context) ([]Permission, error) {
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

func (r *repository) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]Permission, error) {
	var permissions []Permission

	const query = `
		SELECT
			p.id AS id,
			p.code AS code
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	if err := r.db.SelectContext(ctx, &permissions, query, roleID); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *repository) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]Permission, error) {
	var permissions []Permission

	const query = `
		SELECT DISTINCT
			p.id,
			p.code
		FROM users u
		JOIN role_permissions rp ON rp.role_id = u.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1
	`

	if err := r.db.SelectContext(ctx, &permissions, query, userID); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *repository) GetUserRoleCode(ctx context.Context, userID uuid.UUID) (string, error) {
	var code string

	const query = `
		SELECT r.code
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	if err := r.db.GetContext(ctx, &code, query, userID); err != nil {
		return "", err
	}

	return code, nil
}

func (r *repository) SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}

	if len(permissionIDs) > 0 {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, UNNEST($2::bigint[])
		`, roleID, pq.Array(permissionIDs)); err != nil {
			if database.IsForeignKeyViolation(err) {
				return ErrPermissionNotFound
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) GetUserIDsByRoleID(ctx context.Context, roleID int64) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	const query = `SELECT id FROM users WHERE role_id = $1`

	if err := r.db.SelectContext(ctx, &ids, query, roleID); err != nil {
		return nil, err
	}

	return ids, nil
}
