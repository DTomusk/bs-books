package extraction

import "context"

type externalBookModel struct {
	Title   string
	Authors []string
}

type BooksProvider interface {
	SearchBooks(query string, maxResults int, ctx context.Context) ([]externalBookModel, error)
}
