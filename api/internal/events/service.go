package events

import (
	"bs-books-api/internal/db"
	"context"
)

type EventService struct {
	txRunner    db.TxRunner
	repo        *eventRepo
	maxAttempts int
}

// Trev was here
func NewEventService(txRunner db.TxRunner, repo *eventRepo, maxAttempts int) *EventService {
	return &EventService{
		txRunner:    txRunner,
		repo:        repo,
		maxAttempts: maxAttempts,
	}
}

// Write event to queue
func (s *EventService) PublishEvent(ctx context.Context, eventType, aggregateID string, payload any) error {
	event, err := newEvent(eventType, aggregateID, payload)
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
func (s *EventService) DequeueEvent(ctx context.Context) (*Event, error) {
	event, err := s.repo.dequeueEvent(ctx, s.txRunner.DB(), s.maxAttempts)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) MarkEventProcessed(ctx context.Context, db db.DBTX, eventID string) error {
	return s.repo.markEventProcessed(ctx, db, eventID)
}
