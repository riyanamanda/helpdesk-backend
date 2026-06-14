package notification

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/firebase"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type DeviceRepository interface {
	GetTokensByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]string, error)
}

type Notifier interface {
	NewTicket(ctx context.Context, ticketID int64, submitterID uuid.UUID)
	TicketAssigned(ctx context.Context, ticketID int64, assigneeID uuid.UUID, actorID uuid.UUID)
	TicketInProgress(ctx context.Context, ticketID int64, submitterID uuid.UUID, actorID uuid.UUID)
	TicketClosed(ctx context.Context, ticketID int64, createdByID uuid.UUID, actorID uuid.UUID)
	FeedbackStatusUpdated(ctx context.Context, feedbackID int64, createdByID uuid.UUID, actorID uuid.UUID, status string)
}

type notifier struct {
	repo       NotificationRepository
	userRepo   user.UserRepository
	deviceRepo DeviceRepository
	fcm        firebase.FCMSender
}

func NewNotifier(repo NotificationRepository, userRepo user.UserRepository, deviceRepo DeviceRepository, fcm firebase.FCMSender) Notifier {
	return &notifier{
		repo:       repo,
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		fcm:        fcm,
	}
}

func (n *notifier) NewTicket(ctx context.Context, ticketID int64, submitterID uuid.UUID) {
	go func() {
		ctx := context.WithoutCancel(ctx)

		submitterName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, submitterID); err == nil {
			submitterName = u.Name
		}

		adminIDs, err := n.userRepo.GetIDsByRoleAndDivision(ctx, rbac.ADMIN, "IT")
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
			return
		}

		tokens, err := n.deviceRepo.GetTokensByUserIDs(ctx, adminIDs)
		if err != nil {
			slog.WarnContext(ctx, "fcm: failed to get tokens", "error", err)
			return
		}

		if err := n.fcm.SendMulticast(ctx, tokens, buildPayload(NewTicket, TicketReferenceType, ticketID, submitterName, "")); err != nil {
			slog.WarnContext(ctx, "fcm: send failed", "error", err)
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

func (n *notifier) TicketInProgress(ctx context.Context, ticketID int64, submitterID uuid.UUID, actorID uuid.UUID) {
	if submitterID == actorID {
		return
	}

	go func() {
		ctx := context.WithoutCancel(ctx)

		actorName := "Unknown"
		if u, err := n.userRepo.GetByID(ctx, actorID); err == nil {
			actorName = u.Name
		}

		n.notifySingle(ctx, submitterID, TicketInProgress, TicketReferenceType, ticketID, NotificationMetadata{ActorName: actorName})
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
		return
	}

	tokens, err := n.deviceRepo.GetTokensByUserID(ctx, recipientID)
	if err != nil {
		slog.WarnContext(ctx, "fcm: failed to get tokens", "error", err)
		return
	}

	if err := n.fcm.SendMulticast(ctx, tokens, buildPayload(nType, refType, refID, meta.ActorName, meta.Status)); err != nil {
		slog.WarnContext(ctx, "fcm: send failed", "error", err)
	}
}

func buildPayload(nType NotificationType, refType NotificationReferenceType, refID int64, actorName, status string) map[string]string {
	data := map[string]string{
		"type":           string(nType),
		"actor_name":     actorName,
		"reference_type": string(refType),
		"reference_id":   strconv.FormatInt(refID, 10),
	}
	if status != "" {
		data["status"] = status
	}
	return data
}
