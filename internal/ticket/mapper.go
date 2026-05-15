package ticket

func toTicketResponse(t Ticket) TicketResponse {
	return TicketResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		CategoryID:  t.CategoryID,
		Status:      TicketStatus(t.Status),
		Priority:    (*TicketPriority)(t.Priority),
		CreatedBy:   t.CreatedBy,
		AssignedTo:  t.AssignedTo,
		AssignedAt:  t.AssignedAt,
		ResolvedAt:  t.ResolvedAt,
		ClosedAt:    t.ClosedAt,
		ClosedBy:    t.ClosedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toTicketResponses(tickets []Ticket) []TicketResponse {
	responses := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		responses[i] = toTicketResponse(t)
	}
	return responses
}
