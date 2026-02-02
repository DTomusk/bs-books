package ratings

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/reviews"
	"context"
	"database/sql"
)

type RatingService struct {
	txRunner      db.TxRunner
	repo          *ratingRepo
	bookService   *books.BooksService
	reviewService *reviews.ReviewService
}

func NewRatingService(txRunner db.TxRunner, r *ratingRepo, bs *books.BooksService, rs *reviews.ReviewService) *RatingService {
	return &RatingService{
		txRunner:      txRunner,
		repo:          r,
		bookService:   bs,
		reviewService: rs,
	}
}

func (s *RatingService) CreateRating(bookID string, userID string, heartScore float64, pooScore float64, review string, ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Creating rating", "bookID", bookID, "userID", userID, "heartScore", heartScore, "pooScore", pooScore)
	rating, err := newRating(bookID, userID, heartScore, pooScore)

	if err != nil {
		logger.Error("Failed to create rating object", "error", err)
		return err
	}

	// ensure book exists
	exists, err := s.bookService.BookExists(ctx, bookID)

	// TODO: consider if we want to separate errors
	if err != nil || !exists {
		logger.Error("Book not found or error checking book existence", "error", err, "exists", exists)
		return ErrBookNotFound
	}

	// Ensure user hasn't rated this book before
	existingRating, err := s.repo.getRatingByUserAndBook(userID, bookID, ctx, s.txRunner.DB())
	if err != nil {
		logger.Error("Failed to check for existing rating", "error", err)
		return err
	}

	if existingRating != nil {
		logger.Info("User has already rated this book", "userID", userID, "bookID", bookID)
		return ErrRatingAlreadyExists
	}

	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		err = s.repo.create(rating, ctx, tx)
		if err != nil {
			logger.Error("Failed to create rating", "error", err)
			return err
		}

		if review == "" {
			return nil
		}

		err = s.reviewService.CreateReview(bookID, userID, rating.ID, review, ctx, tx)

		if err != nil {
			logger.Error("Failed to create review", "error", err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// queue background task to update book rating stats

	return nil
}
