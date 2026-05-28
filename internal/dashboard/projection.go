package dashboard

import "time"

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
