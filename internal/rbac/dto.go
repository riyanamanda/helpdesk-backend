package rbac

type RoleResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

type PermissionResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}
