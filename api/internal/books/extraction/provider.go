package extraction

import "context"

type externalBookModel struct {
	Title    string
	Authors  []string
	ImageURL string
	Synopsis string
}

type BooksProvider interface {
	SearchBooks(query string, maxResults int, ctx context.Context) ([]externalBookModel, error)
}
