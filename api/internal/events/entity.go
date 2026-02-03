package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          string
	Type        string
	AggregateID string
	Payload     json.RawMessage
	OccurredAt  time.Time
}

func newEvent(eventType string, payload any) (*Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Event{
		ID:         uuid.NewString(),
		Type:       eventType,
		Payload:    json.RawMessage(data),
		OccurredAt: time.Now().UTC(),
	}, nil
}
