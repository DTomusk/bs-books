package reviews

import (
	"bs-books-api/internal/db"
	"context"
)

type ReviewService struct {
	repo                     *ReviewRepo
	db                       db.DBTX
	reviewVisiblityThreshold int
}

func NewReviewService(repo *ReviewRepo, db db.DBTX, reviewVisiblityThreshold int) *ReviewService {
	return &ReviewService{
		repo:                     repo,
		db:                       db,
		reviewVisiblityThreshold: reviewVisiblityThreshold,
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

func (s *ReviewService) HandleReviewReported(ctx context.Context, tx db.DBTX, reviewID string) error {
	review, err := s.repo.getByID(ctx, tx, reviewID)
	if err != nil {
		return err
	}
	if review == nil {
		return nil // If review doesn't exist, we can ignore the event
	}
	err = s.repo.incrementReportCount(ctx, tx, reviewID, s.reviewVisiblityThreshold)
	if err != nil {
		return err
	}
	return nil
}
