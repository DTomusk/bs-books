package books

import "context"

type externalBookModel struct {
	Title   string
	Authors []string
}

type BooksProvider interface {
	SearchBooks(query string, ctx context.Context) ([]externalBookModel, error)
}
