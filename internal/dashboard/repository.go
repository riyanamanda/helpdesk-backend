package dashboard

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DashboardRepository interface {
	GetSummary(ctx context.Context) (SummaryProjection, error)
	GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendProjection, error)
	GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsProjection, error)
	GetAgentWorkload(ctx context.Context) ([]AgentWorkloadProjection, error)
}

type repository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) DashboardRepository {
	return &repository{db: db}
}

func (r *repository) GetSummary(ctx context.Context) (SummaryProjection, error) {
	var (
		statusRow   StatusStatsProjection
		priorityRow PriorityStatsProjection
	)

	const statusQuery = `
		SELECT
			COUNT(*) FILTER (WHERE status = 'IN_PROGRESS')                                            AS in_progress,
			COUNT(*) FILTER (WHERE status = 'RESOLVED')                                               AS resolved,
			COUNT(*) FILTER (WHERE status = 'CLOSED')                                                 AS closed,
			COUNT(*)                                                                                   AS total,
			COUNT(*) FILTER (WHERE status = 'OPEN' AND assigned_to IS NULL)                           AS unassigned
		FROM tickets
	`

	if err := r.db.GetContext(ctx, &statusRow, statusQuery); err != nil {
		return SummaryProjection{}, err
	}

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
		return SummaryProjection{}, err
	}

	return SummaryProjection{
		Status:   statusRow,
		Priority: priorityRow,
	}, nil
}

func (r *repository) GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendProjection, error) {
	var rows []MonthlyTrendProjection

	const query = `
		SELECT
			EXTRACT(MONTH FROM created_at)::int                          AS month,
			COUNT(*)                                                      AS submitted,
			COUNT(*) FILTER (WHERE status IN ('RESOLVED', 'CLOSED'))     AS resolved,
			COUNT(*) FILTER (WHERE status = 'CLOSED')                    AS closed
		FROM tickets
		WHERE EXTRACT(YEAR FROM created_at) = $1
		GROUP BY month
		ORDER BY month
	`

	if err := r.db.SelectContext(ctx, &rows, query, year); err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repository) GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsProjection, error) {
	var rows []CategoryTicketsProjection

	const query = `
		SELECT
			c.id   AS category_id,
			c.name AS category_name,
			COUNT(*) AS total
		FROM tickets t
		JOIN categories c ON c.id = t.category_id
		WHERE t.status IN ('OPEN', 'IN_PROGRESS')
		GROUP BY c.id, c.name
		ORDER BY total DESC, c.name
		LIMIT 10
	`

	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repository) GetAgentWorkload(ctx context.Context) ([]AgentWorkloadProjection, error) {
	var rows []AgentWorkloadProjection

	const query = `
		WITH agent_in_progress AS (
			SELECT assigned_to AS user_id, COUNT(*) AS in_progress
			FROM tickets
			WHERE status = 'IN_PROGRESS'
			GROUP BY assigned_to
		),
		agent_resolved AS (
			SELECT resolved_by AS user_id, COUNT(*) AS resolved
			FROM tickets
			WHERE status IN ('RESOLVED', 'CLOSED')
			  AND DATE_TRUNC('month', resolved_at) = DATE_TRUNC('month', NOW())
			GROUP BY resolved_by
		)
		SELECT
			u.id::text                          AS agent_id,
			u.name                              AS agent_name,
			COALESCE(ip.in_progress, 0)         AS in_progress,
			COALESCE(r.resolved, 0)             AS resolved
		FROM users u
		JOIN (
			SELECT user_id FROM agent_in_progress
			UNION
			SELECT user_id FROM agent_resolved
		) combined ON combined.user_id = u.id
		LEFT JOIN agent_in_progress ip ON ip.user_id = u.id
		LEFT JOIN agent_resolved r ON r.user_id = u.id
		ORDER BY (COALESCE(ip.in_progress, 0) + COALESCE(r.resolved, 0)) DESC, u.name
		LIMIT 10
	`

	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	return rows, nil
}
