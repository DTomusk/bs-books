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

func NewContentModerationService(db db.DBTX, repo *ContentModerationRepo, eventService *events.EventService, reviewService *reviews.ReviewService) *ContentModerationService {
	return &ContentModerationService{
		db:            db,
		repo:          repo,
		eventService:  eventService,
		reviewService: reviewService,
	}
}

func (s *ContentModerationService) ReportContent(ctx context.Context, contentID, contentType, reason, userID string) error {
	switch contentType {
	case Review:
		existing_review, err := s.reviewService.GetReviewByID(ctx, contentID)
		if err != nil || existing_review == nil {
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
		report := NewContentModerationReport(userID, contentID, contentType, existing_review.Text, reason, StatusPending)

		err = s.repo.CreateReport(ctx, s.db, report)
		if err != nil {
			return err
		}

		// Publish event to update review status
		err = s.eventService.PublishEvent(ctx,
			EventReviewReported,
			contentID,
			ReviewReportedEventPayload{},
		)
	default:
		return ErrInvalidContentType
	}

	return nil
}
