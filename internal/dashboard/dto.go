package dashboard

type TicketStatusStats struct {
	Open       int64 `json:"open"`
	InProgress int64 `json:"in_progress"`
	Resolved   int64 `json:"resolved"`
	Closed     int64 `json:"closed"`
	Total      int64 `json:"total"`
}

type TicketPriorityStats struct {
	Low    int64 `json:"low"`
	Medium int64 `json:"medium"`
	High   int64 `json:"high"`
	Urgent int64 `json:"urgent"`
}

type SummaryResponse struct {
	Status   TicketStatusStats   `json:"status"`
	Priority TicketPriorityStats `json:"priority"`
}

type RecentTicketResponse struct {
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	Status     string  `json:"status"`
	Priority   *string `json:"priority"`
	CreatedBy  string  `json:"created_by"`
	AssignedTo *string `json:"assigned_to"`
	CreatedAt  string  `json:"created_at"`
}
