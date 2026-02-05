package content_moderation

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/events"
	"bs-books-api/internal/reviews"
	"context"
)

type ContentModerationService struct {
	db            db.DBTX
	reviewService *reviews.ReviewService
	repo          *ContentModerationRepo
	eventService  *events.EventService
}

func NewContentModerationService(db db.DBTX, reviewService *reviews.ReviewService, repo *ContentModerationRepo, eventService *events.EventService) *ContentModerationService {
	return &ContentModerationService{
		db:            db,
		reviewService: reviewService,
		repo:          repo,
		eventService:  eventService,
	}
}

func (s *ContentModerationService) ReportContent(ctx context.Context, contentID, contentType, reason, userID string) error {
	switch contentType {
	case Review:
		exists, err := s.reviewService.GetReviewExists(ctx, contentID)
		if err != nil || !exists {
			return ErrContentElementDoesntExist
		}

		existing_report, err := s.repo.GetReportByUserByContentID(ctx, s.db, userID, contentID)
		if err != nil {
			return err
		}
		if existing_report != nil {
			return ErrAlreadyReported
		}

		// Add report to repo
		err = s.repo.CreateReport(ctx, s.db, userID, contentID, contentType, reason)
		if err != nil {
			return err
		}

		// Publish event to update review status
		err = s.eventService.PublishEvent(ctx,
			EventReviewReported,
			contentID,
			ReviewReportedEventPayload{},
		)
	}
	return nil
}
