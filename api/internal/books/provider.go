package books

type externalBookModel struct {
	Title  string
	Author string
}

type BooksProvider interface {
	SearchBooks(query string) ([]*externalBookModel, error)
}
