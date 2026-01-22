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
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	AuthorIDs []string `json:"author_ids"`
}

func (r *BookReader) GetAllBooksQuery(ctx context.Context) ([]*BookResponse, error) {
	query := `
		SELECT
			b.id,
			b.title,
			ba.author_id
		FROM books b
		LEFT JOIN book_authors ba ON ba.book_id = b.id
		ORDER BY b.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	booksByID := make(map[string]*BookResponse)
	var orderedBooks []*BookResponse

	for rows.Next() {
		var (
			bookID   string
			title    string
			authorID *string
		)

		if err := rows.Scan(&bookID, &title, &authorID); err != nil {
			return nil, err
		}

		book, exists := booksByID[bookID]
		if !exists {
			book = &BookResponse{
				ID:        bookID,
				Title:     title,
				AuthorIDs: []string{},
			}
			booksByID[bookID] = book
			orderedBooks = append(orderedBooks, book)
		}

		if authorID != nil {
			book.AuthorIDs = append(book.AuthorIDs, *authorID)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderedBooks, nil
}
