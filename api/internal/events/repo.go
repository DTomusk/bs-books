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
	query := `SELECT id, event_type, aggregate_id, payload, created_at, processed_at, attempts, last_error
			  FROM outbox_events
			  WHERE processed_at IS NULL
			  AND attempts < $1
			  ORDER BY created_at
			  LIMIT 1
			  FOR UPDATE SKIP LOCKED`
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
