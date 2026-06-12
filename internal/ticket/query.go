package ticket

import "fmt"

var ticketSortableColumns = map[string]string{
	"created_at": "t.created_at", "updated_at": "t.updated_at",
	"status": "t.status", "priority": "t.priority",
}

const ticketSelectBase = `
	SELECT
		t.id,
		t.title,
		t.description,
		c.id   AS category_id,
		c.name AS category_name,
		d.id   AS division_id,
		d.name AS division_name,
		t.status,
		t.priority,
		u.id   AS created_by_id,
		u.name AS created_by_name,
		uat.id   AS assigned_to_id,
		uat.name AS assigned_to_name,
		uab.id   AS assigned_by_id,
		uab.name AS assigned_by_name,
		urb.id   AS resolved_by_id,
		urb.name AS resolved_by_name,
		ucb.id   AS closed_by_id,
		ucb.name AS closed_by_name,
		t.resolution,
		t.assign_note,
		t.assigned_at,
		t.resolved_at,
		t.closed_at,
		t.created_at,
		t.updated_at
	FROM tickets t
	JOIN categories c ON c.id = t.category_id
	JOIN divisions d ON d.id = t.division_id
	JOIN users u ON u.id = t.created_by
	LEFT JOIN users uat ON uat.id = t.assigned_to
	LEFT JOIN users uab ON uab.id = t.assigned_by
	LEFT JOIN users urb ON urb.id = t.resolved_by
	LEFT JOIN users ucb ON ucb.id = t.closed_by
`

func buildTicketWhere(params GetTicketParams) (string, []any) {
	var (
		where = "WHERE 1=1"
		args  []any
	)

	if params.Status != "" {
		args = append(args, params.Status)
		where += fmt.Sprintf(" AND t.status = $%d", len(args))
	}

	if params.Priority != "" {
		args = append(args, params.Priority)
		where += fmt.Sprintf(" AND t.priority = $%d", len(args))
	}

	if params.CategoryID != nil {
		args = append(args, *params.CategoryID)
		where += fmt.Sprintf(" AND t.category_id = $%d", len(args))
	}

	if params.DivisionID != nil {
		args = append(args, *params.DivisionID)
		where += fmt.Sprintf(" AND t.division_id = $%d", len(args))
	}

	if params.AssignedToID != nil {
		args = append(args, *params.AssignedToID)
		where += fmt.Sprintf(" AND t.assigned_to = $%d", len(args))
	}

	return where, args
}

func buildTicketSort(params GetTicketParams) (string, string) {
	col, ok := ticketSortableColumns[params.SortBy]
	if !ok {
		col = "t.created_at"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return col, dir
}
