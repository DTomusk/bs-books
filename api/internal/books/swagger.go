package books

import "bs-books-api/internal/queries"

// BookListResponse godoc
// @Description Response containing a list of books
type BookListResponse struct {
	Data []*queries.BookResponse `json:"data"`
}
