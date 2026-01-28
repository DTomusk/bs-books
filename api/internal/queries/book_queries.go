package queries

import (
	"bs-books-api/internal/db"
	"context"
)

type BookReader struct {
	db db.DBTX
}

type AuthorSearchItem struct {
	ID   string
	Name string
}

type BookSearchItem struct {
	ID         string
	Title      string
	Similarity float64
	ImageURL   string
	Synopsis   string
	Authors    []AuthorSearchItem
}

type BookSearchPage struct {
	Items      []BookSearchItem
	Total      int
	TotalPages int
	Page       int
	Size       int
}

func NewBookReader(db db.DBTX) *BookReader {
	return &BookReader{db: db}
}

func (r *BookReader) SearchBooksQuery(queryStr string, page, pageSize, offset int, ctx context.Context) (*BookSearchPage, error) {
	// Search books
	const booksQuery = `
	WITH ranked_books AS (
    SELECT
        b.id,
        b.title,
		b.cover_img_url,
		b.synopsis,
        similarity(b.title, $1) AS score
    FROM books b
    WHERE b.title % $1
    ORDER BY score DESC
    LIMIT $2 OFFSET $3
)
	SELECT
	rb.id AS book_id,
	rb.title AS book_title,
	a.id AS author_id,
	a.name AS author_name,
	rb.score AS book_score,
	rb.cover_img_url AS book_image_url,
	rb.synopsis AS book_synopsis
FROM ranked_books rb
LEFT JOIN book_author ba ON rb.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
ORDER BY rb.score DESC, rb.title ASC;
	`

	rows, err := r.db.QueryContext(ctx, booksQuery, queryStr, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookMap := make(map[string]*BookSearchItem)
	order := make([]string, 0)

	for rows.Next() {
		var bookID, bookTitle, authorID, authorName, bookImageURL, bookSynopsis string
		var bookScore float64
		if err := rows.Scan(&bookID, &bookTitle, &authorID, &authorName, &bookScore, &bookImageURL, &bookSynopsis); err != nil {
			return nil, err
		}

		book, exists := bookMap[bookID]
		if !exists {
			book = &BookSearchItem{
				ID:         bookID,
				Title:      bookTitle,
				Similarity: bookScore,
				ImageURL:   bookImageURL,
				Synopsis:   bookSynopsis,
				Authors:    []AuthorSearchItem{},
			}
			bookMap[bookID] = book
			order = append(order, bookID)
		}

		book.Authors = append(book.Authors, AuthorSearchItem{
			ID:   authorID,
			Name: authorName,
		})
	}

	items := make([]BookSearchItem, 0, len(order))
	for _, bookID := range order {
		items = append(items, *bookMap[bookID])
	}

	// Get total count
	const countQuery = `
	SELECT COUNT(*) FROM books WHERE title % $1;`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, queryStr).Scan(&total); err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &BookSearchPage{
		Items:      items,
		Total:      total,
		TotalPages: totalPages,
		Page:       page,
		Size:       pageSize,
	}, nil
}
