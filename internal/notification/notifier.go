package notification

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type Notifier interface {
	NewTicket(ctx context.Context, ticketID int64, submitterID uuid.UUID)
	TicketAssigned(ctx context.Context, ticketID int64, assigneeID uuid.UUID, actorID uuid.UUID)
	TicketClosed(ctx context.Context, ticketID int64, createdByID uuid.UUID, actorID uuid.UUID)
	FeedbackStatusUpdated(ctx context.Context, feedbackID int64, createdByID uuid.UUID, actorID uuid.UUID, status string)
}

type notifier struct {
	repo     NotificationRepository
	userRepo user.UserRepository
}

func NewNotifier(repo NotificationRepository, userRepo user.UserRepository) Notifier {
	return &notifier{repo: repo, userRepo: userRepo}
}

func (n *notifier) NewTicket(ctx context.Context, ticketID int64, submitterID uuid.UUID) {
	go func() {
		ctx := context.WithoutCancel(ctx)

		submitterName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, submitterID); err == nil {
			submitterName = u.Name
		}

		adminIDs, err := n.userRepo.GetIDsByRoleAndDivision(ctx, user.ADMIN, "IT")
		if err != nil || len(adminIDs) == 0 {
			return
		}

		metadata, err := json.Marshal(NotificationMetadata{ActorName: submitterName})
		if err != nil {
			slog.ErrorContext(ctx, "notification: marshal metadata", "error", err)
			return
		}

		notifications := make([]Notification, len(adminIDs))
		for i, id := range adminIDs {
			notifications[i] = Notification{
				UserID:        id,
				Type:          NewTicket,
				ReferenceType: TicketReferenceType,
				ReferenceID:   ticketID,
				Metadata:      string(metadata),
			}
		}

		if err := n.repo.CreateBatch(ctx, notifications); err != nil {
			slog.ErrorContext(ctx, "notification: create batch failed", "ticket_id", ticketID, "error", err)
		}
	}()
}

func (n *notifier) TicketAssigned(ctx context.Context, ticketID int64, assigneeID uuid.UUID, actorID uuid.UUID) {
	go func() {
		ctx := context.WithoutCancel(ctx)

		actorName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, actorID); err == nil {
			actorName = u.Name
		}

		n.notifySingle(ctx, assigneeID, TicketAssigned, TicketReferenceType, ticketID, NotificationMetadata{ActorName: actorName})
	}()
}

func (n *notifier) TicketClosed(ctx context.Context, ticketID int64, createdByID uuid.UUID, actorID uuid.UUID) {
	if createdByID == actorID {
		return
	}

	go func() {
		ctx := context.WithoutCancel(ctx)

		actorName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, actorID); err == nil {
			actorName = u.Name
		}

		n.notifySingle(ctx, createdByID, TicketClosed, TicketReferenceType, ticketID, NotificationMetadata{ActorName: actorName})
	}()
}

func (n *notifier) FeedbackStatusUpdated(ctx context.Context, feedbackID int64, createdByID uuid.UUID, actorID uuid.UUID, status string) {
	if createdByID == actorID {
		return
	}

	go func() {
		ctx := context.WithoutCancel(ctx)

		actorName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, actorID); err == nil {
			actorName = u.Name
		}

		n.notifySingle(ctx, createdByID, FeedbackStatusUpdated, FeedbackReferenceType, feedbackID, NotificationMetadata{ActorName: actorName, Status: status})
	}()
}

func (n *notifier) notifySingle(ctx context.Context, recipientID uuid.UUID, nType NotificationType, refType NotificationReferenceType, refID int64, meta NotificationMetadata) {
	metadata, err := json.Marshal(meta)
	if err != nil {
		slog.ErrorContext(ctx, "notification: marshal metadata", "error", err)
		return
	}

	if err := n.repo.CreateBatch(ctx, []Notification{{
		UserID:        recipientID,
		Type:          nType,
		ReferenceType: refType,
		ReferenceID:   refID,
		Metadata:      string(metadata),
	}}); err != nil {
		slog.ErrorContext(ctx, "notification: create failed", "type", nType, "ref_id", refID, "error", err)
	}
}
