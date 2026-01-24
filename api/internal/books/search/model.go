package search

type BookResult struct {
	ID      string
	Title   string
	Authors []AuthorResult
}

type AuthorResult struct {
	ID   string
	Name string
}

type Page struct {
	Items      []BookResult
	Page       int
	Size       int
	TotalItems int
	TotalPages int
}
