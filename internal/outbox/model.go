package outbox

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func NewEvent(eventType string, aggregateID int64, payload any) (OutboxEvent, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return OutboxEvent{}, fmt.Errorf("encode %s event: %w", eventType, err)
	}

	return OutboxEvent{
		EventType:   eventType,
		AggregateID: strconv.FormatInt(aggregateID, 10),
		Payload:     data,
	}, nil
}

type OutboxEvent struct {
	ID          uuid.UUID       `db:"id"`
	EventType   string          `db:"event_type"`
	AggregateID string          `db:"aggregate_id"`
	Payload     json.RawMessage `db:"payload"`
	CreatedAt   time.Time       `db:"created_at"`
	PublishedAt *time.Time      `db:"published_at"`
}
