package dashboard

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type DashboardRepository interface {
	GetSummary(ctx context.Context) (SummaryResponse, error)
	GetRecentTickets(ctx context.Context) ([]RecentTicket, error)
}

type repository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) DashboardRepository {
	return &repository{db: db}
}

type statusStatsRow struct {
	Open       int64 `db:"open"`
	InProgress int64 `db:"in_progress"`
	Resolved   int64 `db:"resolved"`
	Closed     int64 `db:"closed"`
	Total      int64 `db:"total"`
}

type priorityStatsRow struct {
	Low    int64 `db:"low"`
	Medium int64 `db:"medium"`
	High   int64 `db:"high"`
	Urgent int64 `db:"urgent"`
}

type recentTicketRow struct {
	ID         int64     `db:"id"`
	Title      string    `db:"title"`
	Status     string    `db:"status"`
	Priority   *string   `db:"priority"`
	CreatedBy  string    `db:"created_by_name"`
	AssignedTo *string   `db:"assigned_to_name"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *repository) GetSummary(ctx context.Context) (SummaryResponse, error) {
	var statusRow statusStatsRow
	const statusQuery = `
		SELECT
			COUNT(*) FILTER (WHERE status = 'OPEN')        AS open,
			COUNT(*) FILTER (WHERE status = 'IN_PROGRESS') AS in_progress,
			COUNT(*) FILTER (WHERE status = 'RESOLVED')    AS resolved,
			COUNT(*) FILTER (WHERE status = 'CLOSED')      AS closed,
			COUNT(*)                                        AS total
		FROM tickets
	`
	if err := r.db.GetContext(ctx, &statusRow, statusQuery); err != nil {
		return SummaryResponse{}, err
	}

	var priorityRow priorityStatsRow
	const priorityQuery = `
		SELECT
			COUNT(*) FILTER (WHERE priority = 'LOW')    AS low,
			COUNT(*) FILTER (WHERE priority = 'MEDIUM') AS medium,
			COUNT(*) FILTER (WHERE priority = 'HIGH')   AS high,
			COUNT(*) FILTER (WHERE priority = 'URGENT') AS urgent
		FROM tickets
		WHERE status IN ('OPEN', 'IN_PROGRESS')
	`
	if err := r.db.GetContext(ctx, &priorityRow, priorityQuery); err != nil {
		return SummaryResponse{}, err
	}

	return SummaryResponse{
		ByStatus: TicketStatusStats{
			Open:       statusRow.Open,
			InProgress: statusRow.InProgress,
			Resolved:   statusRow.Resolved,
			Closed:     statusRow.Closed,
			Total:      statusRow.Total,
		},
		ByPriority: TicketPriorityStats{
			Low:    priorityRow.Low,
			Medium: priorityRow.Medium,
			High:   priorityRow.High,
			Urgent: priorityRow.Urgent,
		},
	}, nil
}

func (r *repository) GetRecentTickets(ctx context.Context) ([]RecentTicket, error) {
	var rows []recentTicketRow
	const query = `
		SELECT
			t.id,
			t.title,
			t.status,
			t.priority,
			u.name   AS created_by_name,
			uat.name AS assigned_to_name,
			t.created_at
		FROM tickets t
		JOIN  users u   ON u.id  = t.created_by
		LEFT JOIN users uat ON uat.id = t.assigned_to
		WHERE t.status IN ('OPEN', 'IN_PROGRESS')
		ORDER BY t.created_at DESC
		LIMIT 5
	`
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	tickets := make([]RecentTicket, len(rows))
	for i, row := range rows {
		tickets[i] = RecentTicket{
			ID:         row.ID,
			Title:      row.Title,
			Status:     row.Status,
			Priority:   row.Priority,
			CreatedBy:  row.CreatedBy,
			AssignedTo: row.AssignedTo,
			CreatedAt:  row.CreatedAt.Format(time.RFC3339),
		}
	}

	return tickets, nil
}
