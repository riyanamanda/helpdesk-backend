package ticket

import (
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func toTicketResponse(t TicketProjection) TicketResponse {
	var (
		assignedTo *user.UserBrief
		closedBy   *user.UserBrief
	)

	if t.AssignedToID != nil && t.AssignedToName != nil {
		assignedTo = &user.UserBrief{
			ID:   *t.AssignedToID,
			Name: *t.AssignedToName,
		}
	}

	if t.ClosedByID != nil && t.ClosedByName != nil {
		closedBy = &user.UserBrief{
			ID:   *t.ClosedByID,
			Name: *t.ClosedByName,
		}
	}

	return TicketResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Category: category.CategoryBrief{
			ID:   int64(t.CategoryID),
			Name: t.CategoryName,
		},
		Status:   TicketStatus(t.Status),
		Priority: (*TicketPriority)(t.Priority),
		CreatedBy: user.UserBrief{
			ID:   t.CreatedByID,
			Name: string(t.CreatedByName),
		},
		AssignedTo: assignedTo,
		AssignedAt: t.AssignedAt,
		ResolvedAt: t.ResolvedAt,
		ClosedAt:   t.ClosedAt,
		ClosedBy:   closedBy,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func toTicketResponses(tickets []TicketProjection) []TicketResponse {
	responses := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		responses[i] = toTicketResponse(t)
	}
	return responses
}
