package content_moderation

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/events"
	"bs-books-api/internal/logging"
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

	logger := logging.FromContext(ctx)

	logger.Info("Reporting content", "contentID", contentID, "contentType", contentType, "reason", reason, "userID", userID)
	switch contentType {
	case Review:
		logger.Info("Reporting review")
		existing_review, err := s.reviewService.GetReviewByID(ctx, contentID)
		if err != nil || existing_review == nil {
			logger.Error("Review not found", "contentID", contentID, "error", err)
			return ErrContentElementDoesntExist
		}
		contentSnapshot = existing_review.Text
		eventType = EventReviewReported
		eventPayload = ReviewReportedEventPayload{}
	case User:
		logger.Info("Reporting user")
		user, err := s.userService.GetUserByID(contentID, ctx)
		if err != nil || user == nil {
			logger.Error("User not found", "contentID", contentID, "error", err)
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
		logger.Error("Failed to check existing report", "contentID", contentID, "userID", userID, "error", err)
		return err
	}
	if existing_report != nil {
		logger.Warn("User has already reported this content", "contentID", contentID, "userID", userID)
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

	logger.Info("Content reported successfully", "contentID", contentID, "contentType", contentType, "userID", userID)

	return nil
}
