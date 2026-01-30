package books

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"context"
	"database/sql"
)

type BooksService struct {
	txRunner db.TxRunner
	repo     *booksRepo
}

func NewBooksService(txRunner db.TxRunner, repo *booksRepo) *BooksService {
	return &BooksService{
		txRunner: txRunner,
		repo:     repo,
	}
}

func (s *BooksService) CreateBooksWithAuthors(books []*Book, ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Creating multiple books", "count", len(books))
	for _, book := range books {
		err := s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
			return s.createBookWithAuthors(book, tx, ctx)
		})
		if err != nil {
			logger.Error("Failed to create book with authors", "title", book.Title, "error", err)
		}
	}
	return nil
}

// The service shouldn't know about how book authors are stored
// But at the same time, the repo shouldn't know about transactions
// The service coordinates the transaction
// We want the book to fail if an author association fails
func (s *BooksService) createBookWithAuthors(book *Book, tx *sql.Tx, ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Creating book", "title", book.Title, "id", book.ID, "authors", book.AuthorIDs)
	err := s.repo.createBook(book, ctx, tx)

	if err != nil {
		logger.Error("Failed to create book", "title", book.Title, "error", err)
		return err
	}

	err = s.repo.addAuthorsToBook(book.ID, book.AuthorIDs, ctx, tx)
	if err != nil {
		logger.Error("Failed to add authors to book", "title", book.Title, "error", err)
		return err
	}

	return nil
}

func (s *BooksService) BookExists(ctx context.Context, bookID string) (bool, error) {
	return s.repo.getBookExists(bookID, ctx, s.txRunner.DB())
}
