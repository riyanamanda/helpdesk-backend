package ticket

import (
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func toTicketResponse(t TicketProjection) TicketResponse {

	var (
		assignedTo *user.UserBrief

		resolvedBy *user.UserBrief

		closedBy *user.UserBrief
	)

	if t.AssignedToID != nil && t.AssignedToName != nil {

		assignedTo = &user.UserBrief{

			ID: *t.AssignedToID,

			Name: *t.AssignedToName,
		}

	}

	if t.ResolvedByID != nil && t.ResolvedByName != nil {

		resolvedBy = &user.UserBrief{

			ID: *t.ResolvedByID,

			Name: *t.ResolvedByName,
		}

	}

	if t.ClosedByID != nil && t.ClosedByName != nil {

		closedBy = &user.UserBrief{

			ID: *t.ClosedByID,

			Name: *t.ClosedByName,
		}

	}

	return TicketResponse{

		ID: t.ID,

		Title: t.Title,

		Description: t.Description,

		Category: CategoryBrief{

			ID: t.CategoryID,

			Name: t.CategoryName,
		},

		Division: DivisionBrief{

			ID: t.DivisionID,

			Name: t.DivisionName,
		},

		Status: TicketStatus(t.Status),

		Priority: (*TicketPriority)(t.Priority),

		CreatedBy: user.UserBrief{

			ID: t.CreatedByID,

			Name: t.CreatedByName,
		},

		AssignedTo: assignedTo,

		ResolvedBy: resolvedBy,

		Resolution: t.Resolution,

		AssignedAt: t.AssignedAt,

		ResolvedAt: t.ResolvedAt,

		ClosedAt: t.ClosedAt,

		ClosedBy: closedBy,

		CreatedAt: t.CreatedAt,

		UpdatedAt: t.UpdatedAt,
	}

}

func toTicketAttachmentResponse(ta TicketAttachmentProjection, storageConfig config.Storage) TicketAttachmentResponse {

	return TicketAttachmentResponse{

		ID: ta.ID,

		TicketID: ta.TicketID,

		FileURL: httputil.BuildPublicURL(storageConfig.PublicURL, storageConfig.Bucket, ta.FileKey),

		AttachmentType: ta.AttachmentType,

		UploadedBy: user.UserBrief{

			ID: ta.UploadedByID,

			Name: ta.UploadedByName,
		},

		CreatedAt: ta.CreatedAt,
	}

}

func toTicketAttachmentResponses(attachments []TicketAttachmentProjection, storageConfig config.Storage) []TicketAttachmentResponse {

	return sliceutil.Map(attachments, func(a TicketAttachmentProjection) TicketAttachmentResponse {

		return toTicketAttachmentResponse(a, storageConfig)

	})

}

func toTicketResponses(tickets []TicketProjection) []TicketResponse {

	return sliceutil.Map(tickets, toTicketResponse)

}

func toTicketDetailResponse(ticket TicketProjection, attachments *[]TicketAttachmentProjection, storageConfig config.Storage) TicketDetailResponse {

	var attachmentResponses *[]TicketAttachmentResponse

	if attachments != nil {

		mappedAttachment := toTicketAttachmentResponses(*attachments, storageConfig)

		attachmentResponses = &mappedAttachment

	}

	return TicketDetailResponse{

		TicketResponse: toTicketResponse(ticket),

		Attachments: attachmentResponses,
	}

}
