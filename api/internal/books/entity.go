package books

type Book struct {
	ID        string
	Title     string
	AuthorIDs []string
}

func NewBook(title string, authorIDs []string) *Book {
	return &Book{
		Title:     title,
		AuthorIDs: authorIDs,
	}
}
