package dashboard

import "time"

type SummaryProjection struct {
	Status   StatusStatsProjection
	Priority PriorityStatsProjection
}

type StatusStatsProjection struct {
	InProgress int64 `db:"in_progress"`
	Resolved   int64 `db:"resolved"`
	Closed     int64 `db:"closed"`
	Total      int64 `db:"total"`
	Unassigned int64 `db:"unassigned"`
	Stale      int64 `db:"stale"`
}

type PriorityStatsProjection struct {
	Low    int64 `db:"low"`
	Medium int64 `db:"medium"`
	High   int64 `db:"high"`
	Urgent int64 `db:"urgent"`
}

type RecentTicketProjection struct {
	ID         int64     `db:"id"`
	Title      string    `db:"title"`
	Status     string    `db:"status"`
	Priority   *string   `db:"priority"`
	CreatedBy  string    `db:"created_by_name"`
	AssignedTo *string   `db:"assigned_to_name"`
	CreatedAt  time.Time `db:"created_at"`
}

type MonthlyTrendProjection struct {
	Month     int   `db:"month"`
	Submitted int64 `db:"submitted"`
	Resolved  int64 `db:"resolved"`
	Closed    int64 `db:"closed"`
}

type AgentWorkloadProjection struct {
	AgentID    string `db:"agent_id"`
	AgentName  string `db:"agent_name"`
	InProgress int64  `db:"in_progress"`
	Resolved   int64  `db:"resolved"`
}

type CategoryTicketsProjection struct {
	CategoryID   int64  `db:"category_id"`
	CategoryName string `db:"category_name"`
	Total        int64  `db:"total"`
}
