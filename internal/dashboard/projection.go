package dashboard

import "time"

type SummaryProjection struct {
	Status   StatusStatsProjection
	Priority PriorityStatsProjection
}

type StatusStatsProjection struct {
	Open       int64 `db:"open"`
	InProgress int64 `db:"in_progress"`
	Resolved   int64 `db:"resolved"`
	Closed     int64 `db:"closed"`
	Total      int64 `db:"total"`
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
