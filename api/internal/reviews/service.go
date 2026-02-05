package reviews

import (
	"bs-books-api/internal/db"
	"context"
)

type ReviewService struct {
	repo *ReviewRepo
	db   db.DBTX
}

func NewReviewService(repo *ReviewRepo, db db.DBTX) *ReviewService {
	return &ReviewService{
		repo: repo,
		db:   db,
	}
}

// Note: we create this in a transaction with rating creation
func (s *ReviewService) CreateReview(ratingID, reviewText string, ctx context.Context, tx db.DBTX) error {
	review, err := newReview(ratingID, reviewText)

	if err != nil {
		return err
	}

	err = s.repo.create(review, ctx, tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ReviewService) GetReviewExists(ctx context.Context, reviewID string) (bool, error) {
	review, err := s.repo.getByID(ctx, s.db, reviewID)
	if err != nil {
		return false, err
	}
	return review != nil, nil
}
