package search

import (
	"bs-books-api/internal/db"
	"context"
)

type BookSearchRepo struct{}

func NewBookSearchRepo() *BookSearchRepo {
	return &BookSearchRepo{}
}

func (r *BookSearchRepo) QuerySearchedToday(normalisedQuery string, ctx context.Context, db db.DBTX) (bool, error) {
	const query = `
	SELECT EXISTS (
		SELECT 
		FROM external_search_queries
		WHERE query = $1 AND created_at >= NOW() - INTERVAL '24 hours'
	);
	`

	var exists bool
	err := db.QueryRowContext(ctx, query, normalisedQuery).Scan(&exists)

	return exists, err
}

func (r *BookSearchRepo) LogExternalSearchQuery(normalisedQuery string, ctx context.Context, db db.DBTX) error {
	const insertQuery = `
	INSERT INTO external_search_queries (query, created_at)
	VALUES ($1, NOW())
	ON CONFLICT (query) DO UPDATE SET created_at = NOW();
	`
	_, err := db.ExecContext(ctx, insertQuery, normalisedQuery)
	return err
}
