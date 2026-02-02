package reviews

import (
	"bs-books-api/internal/db"
	"context"
)

type ReviewService struct {
	repo *ReviewRepo
}

func NewReviewService(repo *ReviewRepo) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

// Note: we create this in a transaction with rating creation
func (s *ReviewService) CreateReview(bookID, userID, ratingID, reviewText string, ctx context.Context, tx db.DBTX) error {
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
