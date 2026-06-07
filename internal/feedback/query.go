package feedback

import "fmt"

var feedbackSortableColumns = map[string]string{
	"created_at": "f.created_at",
	"updated_at": "f.updated_at",
	"status":     "f.status",
	"type":       "f.type",
}

const feedbackSelectBase = `
	SELECT
		f.id,
		f.title,
		f.description,
		f.type,
		f.status,
		uc.id   AS created_by_id,
		uc.name AS created_by_name,
		ur.id   AS reviewed_by_id,
		ur.name AS reviewed_by_name,
		f.reviewed_at,
		f.created_at,
		f.updated_at
	FROM feedbacks f
	JOIN users uc ON uc.id = f.created_by
	LEFT JOIN users ur ON ur.id = f.reviewed_by
`

func buildFeedbackWhere(params GetFeedbackParams) (string, []any) {
	var (
		where = "WHERE 1=1"
		args  []any
	)

	if params.Type != nil {
		args = append(args, *params.Type)
		where += fmt.Sprintf(" AND f.type = $%d", len(args))
	}

	if params.Status != nil {
		args = append(args, *params.Status)
		where += fmt.Sprintf(" AND f.status = $%d", len(args))
	}

	if params.CreatedByID != nil {
		args = append(args, *params.CreatedByID)
		where += fmt.Sprintf(" AND f.created_by = $%d", len(args))
	}

	return where, args
}

func buildFeedbackSort(params GetFeedbackParams) (string, string) {
	col, ok := feedbackSortableColumns[params.SortBy]
	if !ok {
		col = "f.created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return col, dir
}
