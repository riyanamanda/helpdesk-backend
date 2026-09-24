package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
)

type Publisher struct {
	outboxRepo outbox.Repository
	rabbitmq   *rabbitmq.Client
}

func NewPublisher(outboxRepo outbox.Repository, rabbitmq *rabbitmq.Client) *Publisher {
	return &Publisher{
		outboxRepo: outboxRepo,
		rabbitmq:   rabbitmq,
	}
}

func (p *Publisher) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			p.publish(ctx)
		}
	}
}

func (p *Publisher) publish(ctx context.Context) {
	events, err := p.outboxRepo.GetUnpublished(ctx, 10)
	if err != nil {
		slog.Error("get unpublished events", "error", err)
		return
	}

	for _, event := range events {
		err := p.rabbitmq.Publish(ctx, "helpdesk.events", event.EventType, "application/json", event.Payload)
		if err != nil {
			slog.Error("publish outbox event failed", "id", event.ID, "event_type", event.EventType, "error", err)
			continue
		}

		if err := p.outboxRepo.MarkPublished(ctx, event.ID); err != nil {
			slog.Error("mark outbox event published failed", "id", event.ID, "error", err)
			continue
		}

		slog.Info("publish event successfully", "id", event.ID, "type", event.EventType)
	}
}
