package ratings

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/db"
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
	// validate and create rating object
	rating, err := newRating(bookID, userID, heartScore, pooScore)

	if err != nil {
		return err
	}

	// ensure book exists
	exists, err := s.bookService.BookExists(ctx, bookID)

	// TODO: consider if we want to separate errors
	if err != nil || !exists {
		return ErrBookNotFound
	}

	// Ensure user hasn't rated this book before
	existingRating, err := s.repo.getRatingByUserAndBook(userID, bookID, ctx, s.txRunner.DB())
	if err != nil {
		return err
	}

	if existingRating != nil {
		return ErrRatingAlreadyExists
	}

	// Consider transaction here
	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		err = s.repo.create(rating, ctx, tx)
		if err != nil {
			return err
		}

		if review == "" {
			return nil
		}

		err = s.reviewService.CreateReview(bookID, userID, rating.ID, review, ctx, tx)

		if err != nil {
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
