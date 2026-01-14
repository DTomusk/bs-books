package ratings

import (
	"bs-books-api/internal/db"
	"context"
)

type RatingService struct {
	db   db.DBTX
	repo *ratingRepo
}

func NewRatingService(db db.DBTX, r *ratingRepo) *RatingService {
	return &RatingService{
		db:   db,
		repo: r,
	}
}

func (s *RatingService) CreateRating(bookID string, heartScore float64, pooScore float64, ctx context.Context) (*Rating, error) {
	rating, err := newRating(bookID, heartScore, pooScore)

	if err != nil {
		return nil, err
	}

	err = s.repo.create(rating, ctx, s.db)

	return rating, err
}
