package queries

import "bs-books-api/internal/db"

type BookResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	AuthorID string `json:"author_id"`
}

func GetAllBooksQuery(db db.DBTX) []*BookResponse {
	return []*BookResponse{
		{ID: "1", Title: "Book 1", AuthorID: "1"},
		{ID: "2", Title: "Book 2", AuthorID: "2"},
		{ID: "3", Title: "Book 3", AuthorID: "3"},
	}
}
