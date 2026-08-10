package dashboard

func toSummary(s SummaryProjection) SummaryResponse {
	return SummaryResponse{
		Status:   TicketStatusStats(s.Status),
		Priority: TicketPriorityStats(s.Priority),
	}
}

func toMonthlyTrend(rows []MonthlyTrendProjection) []MonthlyTrendResponse {
	result := make([]MonthlyTrendResponse, len(rows))

	for i, row := range rows {
		result[i] = MonthlyTrendResponse(row)
	}

	return result
}

func toCategoryTickets(rows []CategoryTicketsProjection) []CategoryTicketsResponse {
	result := make([]CategoryTicketsResponse, len(rows))

	for i, row := range rows {
		result[i] = CategoryTicketsResponse{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Total:        row.Total,
		}
	}

	return result
}

func toAgentWorkload(rows []AgentWorkloadProjection) []AgentWorkloadResponse {
	result := make([]AgentWorkloadResponse, len(rows))

	for i, row := range rows {
		result[i] = AgentWorkloadResponse{
			AgentID:    row.AgentID,
			AgentName:  row.AgentName,
			InProgress: row.InProgress,
			Resolved:   row.Resolved,
		}
	}

	return result
}
