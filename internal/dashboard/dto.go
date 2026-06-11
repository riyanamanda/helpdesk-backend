package dashboard

type TicketStatusStats struct {
	InProgress int64 `json:"in_progress"`
	Resolved   int64 `json:"resolved"`
	Closed     int64 `json:"closed"`
	Total      int64 `json:"total"`
	Unassigned int64 `json:"unassigned"`
	Stale      int64 `json:"stale"`
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

type MonthlyTrendResponse struct {
	Month     int   `json:"month"`
	Submitted int64 `json:"submitted"`
	Resolved  int64 `json:"resolved"`
	Closed    int64 `json:"closed"`
}

type AgentWorkloadResponse struct {
	AgentID    string `json:"agent_id"`
	AgentName  string `json:"agent_name"`
	InProgress int64  `json:"in_progress"`
	Resolved   int64  `json:"resolved"`
}

type CategoryTicketsResponse struct {
	CategoryID   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
	Total        int64  `json:"total"`
}
