package user

import "fmt"

var userSortableColumns = map[string]string{
	"name": "u.name", "role": "r.code", "division": "u.division_id",
	"is_active": "u.is_active", "created_at": "u.created_at",
}

const userSelectBase = `
	SELECT
		u.id,
		u.name,
		u.email,
		u.google_id,
		u.avatar_key,
		u.phone,
		u.role_id as role_id,
		r.code as role_name,
		u.gender,
		d.id as division_id,
		d.name as division_name,
		u.is_active,
		cb.id as created_by_id,
		cb.name as created_by_name,
		u.created_at,
		u.updated_at
	FROM users u
	JOIN roles r
		ON r.id = u.role_id
	LEFT JOIN divisions d
		ON d.id = u.division_id
	LEFT JOIN users cb
		ON cb.id = u.created_by
`

const userSelectWithPassword = `
	SELECT
		u.id,
		u.name,
		u.email,
		u.password,
		u.google_id,
		u.avatar_key,
		u.phone,
		u.role_id as role_id,
		r.code as role_name,
		u.gender,
		d.id as division_id,
		d.name as division_name,
		u.is_active,
		cb.id as created_by_id,
		cb.name as created_by_name,
		u.created_at,
		u.updated_at
	FROM users u
	JOIN roles r
		ON r.id = u.role_id
	LEFT JOIN divisions d
		ON d.id = u.division_id
	LEFT JOIN users cb
		ON cb.id = u.created_by
`

func buildUserWhere(params GetUserParams) (string, []any) {
	var (
		where = "WHERE 1=1"
		args  []any
	)

	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		where += fmt.Sprintf(" AND u.name ILIKE $%d", len(args))
	}

	if params.IsActive != nil {
		args = append(args, *params.IsActive)
		where += fmt.Sprintf(" AND u.is_active = $%d", len(args))
	}

	if params.Role.Name != "" {
		args = append(args, params.Role.Name)
		where += fmt.Sprintf(" AND r.code = $%d", len(args))
	}

	if params.Division != nil {
		args = append(args, *params.Division)
		where += fmt.Sprintf(" AND u.division_id = $%d", len(args))
	}

	return where, args
}

func buildUserSort(params GetUserParams) (string, string) {
	col, ok := userSortableColumns[params.SortBy]
	if !ok {
		col = "u.created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return col, dir
}
