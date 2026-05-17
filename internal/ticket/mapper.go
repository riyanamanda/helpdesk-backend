package ticket

import (
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func toTicketResponse(t TicketProjection) TicketResponse {
	var (
		assignedTo *user.UserBrief
		resolvedBy *user.UserBrief
		closedBy   *user.UserBrief
	)

	if t.AssignedToID != nil && t.AssignedToName != nil {
		assignedTo = &user.UserBrief{
			ID:   *t.AssignedToID,
			Name: *t.AssignedToName,
		}
	}

	if t.ResolvedByID != nil && t.ResolvedByName != nil {
		resolvedBy = &user.UserBrief{
			ID:   *t.ResolvedByID,
			Name: *t.ResolvedByName,
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
		ResolvedBy: resolvedBy,
		Resolution: t.Resolution,
		AssignedAt: t.AssignedAt,
		ResolvedAt: t.ResolvedAt,
		ClosedAt:   t.ClosedAt,
		ClosedBy:   closedBy,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func toTicketAttachmentResponse(ta TicketAttachmentProjection, storage storage.Storage) TicketAttachmentResponse {
	return TicketAttachmentResponse{
		ID:             ta.ID,
		TicketID:       ta.TicketID,
		FileURL:        storage.GetURL(ta.FileKey),
		AttachmentType: ta.AttachmentType,
		UploadedBy: user.UserBrief{
			ID:   ta.UploadedByID,
			Name: ta.UploadedByName,
		},
		CreatedAt: ta.CreatedAt,
	}
}

func toTicketResponses(tickets []TicketProjection) []TicketResponse {
	responses := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		responses[i] = toTicketResponse(t)
	}
	return responses
}

func toTicketDetailResponse(ticket TicketProjection, attachment *TicketAttachmentProjection, storageService storage.Storage) TicketDetailResponse {
	var attachmentResponse *TicketAttachmentResponse
	if attachment != nil {
		mappedAttachment := toTicketAttachmentResponse(*attachment, storageService)
		attachmentResponse = &mappedAttachment
	}

	return TicketDetailResponse{
		TicketResponse: toTicketResponse(ticket),
		Attachment:     attachmentResponse,
	}
}
