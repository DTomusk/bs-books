package ratings

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/db"
	"context"
)

type RatingService struct {
	db          db.DBTX
	repo        *ratingRepo
	bookService *books.BooksService
}

func NewRatingService(db db.DBTX, r *ratingRepo, bs *books.BooksService) *RatingService {
	return &RatingService{
		db:          db,
		repo:        r,
		bookService: bs,
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
	existingRating, err := s.repo.getRatingByUserAndBook(userID, bookID, ctx, s.db)
	if err != nil {
		return err
	}

	if existingRating != nil {
		return ErrRatingAlreadyExists
	}

	err = s.repo.create(rating, ctx, s.db)

	if err != nil {
		return err
	}

	// queue background task to update book rating stats

	return nil
}
