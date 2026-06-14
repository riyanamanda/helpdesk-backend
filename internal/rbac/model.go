package rbac

type Role struct {
	ID   int64  `db:"id"`
	Code string `db:"code"`
}

type Permission struct {
	ID   int64  `db:"id"`
	Code string `db:"code"`
}

type RolePermission struct {
	RoleID       int64 `db:"role_id"`
	PermissionID int64 `db:"permission_id"`
}
