DROP INDEX IF EXISTS idx_outbox_events_event_type;

DROP INDEX IF EXISTS idx_outbox_events_aggregate_id;

DROP INDEX IF EXISTS idx_outbox_events_unprocessed;

DROP TABLE IF EXISTS outbox_events;
