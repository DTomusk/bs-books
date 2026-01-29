package queries

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
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

type BookDetails struct {
	ID       string
	Title    string
	ImageURL string
	Synopsis string
	Authors  []AuthorSearchItem
}

func NewBookReader(db db.DBTX) *BookReader {
	return &BookReader{db: db}
}

func (r *BookReader) SearchBooksQuery(normalisedQuery string, page, pageSize, offset int, ctx context.Context) (*BookSearchPage, error) {
	// Search books
	const booksQuery = `
	WITH ranked_books AS (
    SELECT
        b.id,
        b.title,
		b.cover_img_url,
		b.synopsis,
        similarity(b.normalised_title, $1) AS score
    FROM books b
    WHERE b.normalised_title % $1
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

	rows, err := r.db.QueryContext(ctx, booksQuery, normalisedQuery, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookMap := make(map[string]*BookSearchItem)
	order := make([]string, 0)

	for rows.Next() {
		var bookID, bookTitle, authorID, authorName string
		var bookScore float64
		var bookImageURL, bookSynopsis sql.NullString
		if err := rows.Scan(&bookID, &bookTitle, &authorID, &authorName, &bookScore, &bookImageURL, &bookSynopsis); err != nil {
			return nil, err
		}

		book, exists := bookMap[bookID]
		if !exists {
			book = &BookSearchItem{
				ID:         bookID,
				Title:      bookTitle,
				Similarity: bookScore,
				ImageURL:   nullString(bookImageURL),
				Synopsis:   nullString(bookSynopsis),
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
	SELECT COUNT(*) FROM books WHERE normalised_title % $1;`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, normalisedQuery).Scan(&total); err != nil {
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

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (r *BookReader) GetBookByID(ctx context.Context, id string) (*BookDetails, error) {
	const query = `
	SELECT 
		b.id,
		b.title,
		b.cover_img_url,
		b.synopsis,
		a.id AS author_id,
		a.name AS author_name
	FROM books b
	INNER JOIN book_author ba ON b.id = ba.book_id
	INNER JOIN authors a ON ba.author_id = a.id
	WHERE b.id = $1;
	`

	var book BookDetails

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var bookID, bookTitle, authorID, authorName string
		var bookImageURL, bookSynopsis sql.NullString
		if err := rows.Scan(&bookID, &bookTitle, &bookImageURL, &bookSynopsis, &authorID, &authorName); err != nil {
			return nil, err
		}
		// If this is the first iteration, initialize the book details
		if book.ID == "" {
			book.ID = bookID
			book.Title = bookTitle
			book.ImageURL = nullString(bookImageURL)
			book.Synopsis = nullString(bookSynopsis)
			book.Authors = []AuthorSearchItem{}
		}
		book.Authors = append(book.Authors, AuthorSearchItem{
			ID:   authorID,
			Name: authorName,
		})
	}

	if book.ID == "" {
		return nil, nil
	}

	return &book, nil
}
