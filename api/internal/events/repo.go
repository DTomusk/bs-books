package events

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type OutboxEventRow struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     json.RawMessage

	CreatedAt   time.Time
	ProcessedAt sql.NullTime
	Attempts    int
	LastError   sql.NullString
}

func toEntity(row *OutboxEventRow) *Event {
	return &Event{
		ID:          row.ID,
		Type:        row.EventType,
		AggregateID: row.AggregateID,
		Payload:     row.Payload,
		OccurredAt:  row.CreatedAt,
	}
}

type eventRepo struct{}

func NewEventRepo() *eventRepo {
	return &eventRepo{}
}

func (r *eventRepo) insertEvent(ctx context.Context, db db.DBTX, event *Event) error {
	query := `INSERT INTO outbox_events (id, event_type, aggregate_id, payload, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Type,
		event.AggregateID,
		event.Payload,
		event.OccurredAt,
	)
	return err
}

func (r *eventRepo) dequeueEvent(ctx context.Context, db db.DBTX, maxAttempts int) (*Event, error) {
	var row OutboxEventRow
	query := `
	UPDATE outbox_events
	SET 
		attempts = outbox_events.attempts + 1
	FROM (
		SELECT id, event_type, aggregate_id, payload, created_at, processed_at, attempts, last_error
		FROM outbox_events
		WHERE processed_at IS NULL
		AND attempts < $1
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	) event
	WHERE outbox_events.id = event.id
	RETURNING 
		outbox_events.id, 
		outbox_events.event_type, 
		outbox_events.aggregate_id, 
		outbox_events.payload, 
		outbox_events.created_at, 
		outbox_events.processed_at, 
		outbox_events.attempts, 
		outbox_events.last_error;`
	err := db.QueryRowContext(ctx, query, maxAttempts).Scan(
		&row.ID,
		&row.EventType,
		&row.AggregateID,
		&row.Payload,
		&row.CreatedAt,
		&row.ProcessedAt,
		&row.Attempts,
		&row.LastError,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toEntity(&row), nil
}

func (r *eventRepo) markEventProcessed(ctx context.Context, db db.DBTX, eventID string) error {
	query := `
		UPDATE outbox_events
		SET processed_at = $1
		WHERE id = $2;
	`
	_, err := db.ExecContext(ctx, query, time.Now().UTC(), eventID)
	return err
}
