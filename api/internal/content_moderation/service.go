package content_moderation

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/events"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/users"
	"context"
)

type ContentModerationService struct {
	db            db.DBTX
	reviewService *reviews.ReviewService
	repo          *ContentModerationRepo
	eventService  *events.EventService
	userService   *users.UserService
}

func NewContentModerationService(
	db db.DBTX,
	repo *ContentModerationRepo,
	eventService *events.EventService,
	reviewService *reviews.ReviewService,
	userService *users.UserService,
) *ContentModerationService {
	return &ContentModerationService{
		db:            db,
		repo:          repo,
		eventService:  eventService,
		reviewService: reviewService,
		userService:   userService,
	}
}

func (s *ContentModerationService) ReportContent(ctx context.Context, contentID, contentType, reason, userID string) error {
	var contentSnapshot string
	var eventType string
	var eventPayload any

	switch contentType {
	case Review:
		existing_review, err := s.reviewService.GetReviewByID(ctx, contentID)
		if err != nil || existing_review == nil {
			return ErrContentElementDoesntExist
		}
		contentSnapshot = existing_review.Text
		eventType = EventReviewReported
		eventPayload = ReviewReportedEventPayload{}
	case User:
		user, err := s.userService.GetUserByID(contentID, ctx)
		if err != nil || user == nil {
			return ErrContentElementDoesntExist
		}
		contentSnapshot = user.Username
		eventType = EventUserReported
		eventPayload = UserReportedEventPayload{}
	default:
		return ErrInvalidContentType
	}

	existing_report, err := s.repo.GetReportByUserByContentID(ctx, s.db, userID, contentID)
	if err != nil {
		return err
	}
	if existing_report != nil {
		return ErrAlreadyReported
	}

	// Add report to repo
	report := NewContentModerationReport(userID, contentID, contentType, contentSnapshot, reason, StatusPending)

	err = s.repo.CreateReport(ctx, s.db, report)
	if err != nil {
		return err
	}

	// Publish event to update review status
	err = s.eventService.PublishEvent(ctx,
		eventType,
		contentID,
		eventPayload,
	)
	if err != nil {
		return err
	}

	return nil
}
