package books

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
)

type booksRepo struct{}

func NewBooksRepo() *booksRepo {
	return &booksRepo{}
}

func (r *booksRepo) createBook(book *Book, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, `INSERT INTO books (id, title) VALUES ($1, $2)`, book.ID, book.Title)
	return err
}

func (r *booksRepo) addAuthorsToBook(bookID string, authorIDs []string, ctx context.Context, db db.DBTX) error {
	for _, authorID := range authorIDs {
		_, err := db.ExecContext(ctx, `INSERT INTO book_author (book_id, author_id) VALUES ($1, $2)`, bookID, authorID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *booksRepo) getBookByID(bookID string, ctx context.Context, db db.DBTX) (*Book, error) {
	var book Book
	row := db.QueryRowContext(ctx, `SELECT id, title FROM books WHERE id = $1`, bookID)
	err := row.Scan(&book.ID, &book.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	// Get author IDs
	rows, err := db.QueryContext(ctx, `SELECT author_id FROM book_author WHERE book_id = $1`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var authorID string
		if err := rows.Scan(&authorID); err != nil {
			return nil, err
		}
		book.AuthorIDs = append(book.AuthorIDs, authorID)
	}
	return &book, nil
}
