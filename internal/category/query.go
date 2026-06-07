package category

import "fmt"

var categorySortableColumns = map[string]string{
	"name": "name", "is_active": "is_active", "created_at": "created_at",
}

const categorySelectBase = `
	SELECT
		id,
		name,
		is_active,
		created_at,
		updated_at
	FROM categories
`

func buildCategoryWhere(params GetCategoryParams) (string, []any) {
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

func buildCategorySort(params GetCategoryParams) (string, string) {
	col, ok := categorySortableColumns[params.SortBy]
	if !ok {
		col = "created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return col, dir
}
