package division

import "fmt"

var divisionSortableColumns = map[string]string{
	"name": "name", "is_active": "is_active", "created_at": "created_at",
}

const divisionSelectBase = `
	SELECT
		id,
		name,
		is_active,
		created_at,
		updated_at
	FROM divisions
`

func buildDivisionWhere(params GetDivisionParams) (string, []any) {
	var (
		where = "WHERE 1=1"
		args  []any
	)

	if params.Search != "" {
		args = append(args, "%"+params.Search+"%")
		where += fmt.Sprintf(" AND name ILIKE $%d", len(args))
	}

	if params.IsActive != nil {
		args = append(args, *params.IsActive)
		where += fmt.Sprintf(" AND is_active = $%d", len(args))
	}

	return where, args
}

func buildDivisionSort(params GetDivisionParams) (string, string) {
	col, ok := divisionSortableColumns[params.SortBy]
	if !ok {
		col = "created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return col, dir
}
