CREATE TABLE external_search_queries (
    query TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_external_search_queries_query_created_at
    ON external_search_queries (query, created_at);