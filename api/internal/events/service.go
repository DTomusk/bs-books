package events

import (
	"bs-books-api/internal/db"
	"context"
)

type EventService struct {
	txRunner db.TxRunner
	repo     *eventRepo
}

// Trev was here
func NewEventService(txRunner db.TxRunner, repo *eventRepo) *EventService {
	return &EventService{
		txRunner: txRunner,
		repo:     repo,
	}
}

// Write event to queue
func (s *EventService) EnqueueEvent(ctx context.Context, eventType string, payload interface{}) error {
	event, err := newEvent(eventType, payload)
	if err != nil {
		return err
	}

	err = s.repo.insertEvent(ctx, s.txRunner.DB(), event)

	if err != nil {
		return err
	}
	return nil
}

// Get next event to process from queue
func (s *EventService) DequeueEvent(ctx context.Context) error {
	return nil
}
