package rbac

func toRoleResponse(r Role) RoleResponse {
	return RoleResponse(r)
}

func toPermissionResponse(p Permission) PermissionResponse {
	return PermissionResponse(p)
}

func toRoleResponses(roles []Role) []RoleResponse {
	result := make([]RoleResponse, len(roles))
	for i, r := range roles {
		result[i] = toRoleResponse(r)
	}

	return result
}

func toPermissionResponses(permissions []Permission) []PermissionResponse {
	result := make([]PermissionResponse, len(permissions))
	for i, p := range permissions {
		result[i] = toPermissionResponse(p)
	}

	return result
}
