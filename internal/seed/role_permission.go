package seed

import "github.com/jmoiron/sqlx"

func SeedRolePermission(db *sqlx.DB) (int64, error) {
	var rolePermission []struct {
		RoleID       int64
		PermissionID int64
	}

	const queryPermission = `
    SELECT id
    FROM permissions
`

	var permissionIds []int64
	err := db.Select(&permissionIds, queryPermission)
	if err != nil {
		return 0, err
	}

	type RolePermissionStruct struct {
		RoleID       int64
		PermissionID int64
	}

	for _, p := range permissionIds {
		rolePermission = append(rolePermission, struct {
			RoleID       int64
			PermissionID int64
		}{
			RoleID:       1,
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
