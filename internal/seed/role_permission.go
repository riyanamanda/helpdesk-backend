package seed

import "github.com/jmoiron/sqlx"

func SeedRolePermission(db *sqlx.DB) (int64, error) {
	var superadminRoleID int64
	if err := db.Get(&superadminRoleID, `SELECT id FROM roles WHERE code = 'SUPERADMIN'`); err != nil {
		return 0, err
	}

	var permissionIds []int64
	if err := db.Select(&permissionIds, `SELECT id FROM permissions`); err != nil {
		return 0, err
	}

	var rolePermission []struct {
		RoleID       int64
		PermissionID int64
	}
	for _, p := range permissionIds {
		rolePermission = append(rolePermission, struct {
			RoleID       int64
			PermissionID int64
		}{
			RoleID:       superadminRoleID,
			PermissionID: p,
		})
	}

	const query = `
		INSERT INTO role_permissions(role_id, permission_id)
		VALUES($1, $2)
		ON CONFLICT DO NOTHING
	`

	var affected int64
	for _, rp := range rolePermission {
		result, err := db.Exec(query, rp.RoleID, rp.PermissionID)
		if err != nil {
			return 0, err
		}

		n, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		affected += n
	}

	return affected, nil
}
