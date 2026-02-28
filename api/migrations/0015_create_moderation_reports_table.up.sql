CREATE TABLE moderation_reports (
    id UUID PRIMARY KEY,
    
    reporter_id UUID REFERENCES users(id),
    reason TEXT NOT NULL,

    content_type TEXT NOT NULL,
    content_id UUID NOT NULL,
    content_snapshot TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    status TEXT NOT NULL DEFAULT 'pending_review'
);

CREATE UNIQUE INDEX uniq_report_per_user ON moderation_reports (reporter_id, content_type, content_id);