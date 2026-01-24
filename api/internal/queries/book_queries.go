package queries

import (
	"bs-books-api/internal/db"
	"context"
)

type BookReader struct {
	db db.DBTX
}

type BookSearchItem struct {
	ID      string
	Title   string
	Authors []string
}

type BookSearchPage struct {
	Items []BookSearchItem
	Total int
	Page  int
	Size  int
}

func NewBookReader(db db.DBTX) *BookReader {
	return &BookReader{db: db}
}

// Search
func (r *BookReader) SearchBooksQuery(queryStr string, page, pageSize, offset int, ctx context.Context) (*BookSearchPage, error) {
	// Search books
	// Get total count
	return &BookSearchPage{}, nil
}
