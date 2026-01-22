package books

type externalBookModel struct {
	Title   string
	Authors []string
}

type BooksProvider interface {
	SearchBooks(query string) ([]externalBookModel, error)
}
