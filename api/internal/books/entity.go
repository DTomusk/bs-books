package books

import "github.com/google/uuid"

type Book struct {
	ID        string
	Title     string
	AuthorIDs []string
}

func NewBook(title string, authorIDs []string) *Book {
	return &Book{
		ID:        uuid.NewString(),
		Title:     title,
		AuthorIDs: authorIDs,
	}
}
