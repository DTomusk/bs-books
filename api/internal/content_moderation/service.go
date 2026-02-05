package content_moderation

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/reviews"
	"context"
)

type ContentModerationService struct {
	db            db.DBTX
	reviewService *reviews.ReviewService
	repo          *ContentModerationRepo
}

func NewContentModerationService(db db.DBTX, reviewService *reviews.ReviewService, repo *ContentModerationRepo) *ContentModerationService {
	return &ContentModerationService{
		db:            db,
		reviewService: reviewService,
		repo:          repo,
	}
}

func (s *ContentModerationService) ReportContent(ctx context.Context, contentID, contentType, reason, userID string) error {
	// Check that the content exists and is of the correct type (e.g. review)
	switch contentType {
	case Review:
		// call review service to check content existence
		exists, err := s.reviewService.GetReviewByID(contentID)
		if err != nil || !exists {
			return ErrContentElementDoesntExist
		}

		// then check whether the user has reported this content before
		existing_report, err := s.repo.GetReportByUserByContentID(ctx, s.db, userID, contentID)
		if err != nil {
			return err
		}
		if existing_report != nil {
			// return that user has already reported content
		}
	}
	// Check that the user has not already reported this content
	// Create a new report in the database with status "pending_review"
	// Queue event to update review moderation status (this should automatically hide the content if reports exceeds a threshold)
	return nil
}
