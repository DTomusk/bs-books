package queries

import (
	"bs-books-api/internal/db"
	"context"
)

type BookReader struct {
	db db.DBTX
}

func NewBookReader(db db.DBTX) *BookReader {
	return &BookReader{db: db}
}

type BookResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	AuthorID string `json:"author_id"`
}

func (r *BookReader) GetAllBooksQuery(ctx context.Context) ([]*BookResponse, error) {
	var books []*BookResponse
	query := "SELECT id, title, author_id FROM books"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book BookResponse
		if err := rows.Scan(&book.ID, &book.Title, &book.AuthorID); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}
