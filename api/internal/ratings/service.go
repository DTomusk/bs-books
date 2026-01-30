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

func (s *RatingService) CreateRating(bookID string, heartScore float64, pooScore float64, ctx context.Context) error {
	rating, err := newRating(bookID, heartScore, pooScore)

	// ensure book exists
	exists, err := s.bookService.BookExists(ctx, bookID)

	// TODO: consider if we want to separate errors
	if err != nil || !exists {
		return ErrBookNotFound
	}

	err = s.repo.create(rating, ctx, s.db)

	if err != nil {
		return err
	}

	// queue background task to update book rating stats

	return nil
}
