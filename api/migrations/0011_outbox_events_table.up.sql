CREATE TABLE outbox_events (
    id UUID PRIMARY KEY, 
    event_type TEXT NOT NULL,
    aggregate_id UUID NOT NULL,
    payload JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,

    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT
);

CREATE INDEX idx_outbox_events_unprocessed
ON outbox_events (created_at)
WHERE processed_at IS NULL;

CREATE INDEX idx_outbox_events_aggregate_id
ON outbox_events (aggregate_id);

CREATE INDEX idx_outbox_events_event_type
ON outbox_events (event_type);