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

// Update the rating metadata on a book when a rating gets created
func (s *BooksService) AddRatingToBook(ctx context.Context, tx *sql.Tx, bookID string, heartScore float64, pooScore float64) error {
	averageHeartScore, averagePooScore, totalRatings, err := s.repo.getBookRatingStats(bookID, ctx, tx)
	if err != nil {
		return err
	}

	newTotalRatings := totalRatings + 1
	newAverageHeartScore := ((averageHeartScore * float64(totalRatings)) + heartScore) / float64(newTotalRatings)
	newAveragePooScore := ((averagePooScore * float64(totalRatings)) + pooScore) / float64(newTotalRatings)

	err = s.repo.updateBookRatingStats(bookID, newAverageHeartScore, newAveragePooScore, newTotalRatings, ctx, tx)
	if err != nil {
		return err
	}
	return nil
}
