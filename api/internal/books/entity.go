package books

import "github.com/google/uuid"

type Book struct {
	ID              string
	Title           string
	NormalisedTitle string
	AuthorIDs       []string
}

func NewBook(title string, authorIDs []string) (*Book, error) {
	if len(authorIDs) == 0 {
		return nil, ErrNoAuthorsProvided
	}
	return &Book{
		ID:              uuid.NewString(),
		Title:           title,
		NormalisedTitle: NormaliseBookTitle(title),
		AuthorIDs:       authorIDs,
	}, nil
}
