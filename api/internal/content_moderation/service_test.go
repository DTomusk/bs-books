package content_moderation

import (
	"bs-books-api/internal/events"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidContentType_ReturnsError(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		service := NewContentModerationService(tx, nil, nil, nil)
		ctx := context.Background()

		// Act
		err := service.ReportContent(ctx, "some-id", "invalid-type", "reason", "user-id")

		// Assert
		require.NotNil(t, err)
		require.Equal(t, ErrInvalidContentType, err)
	})
}

func TestReportReview_ReviewDoesNotExist_ReturnsError(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), tx, 5)
		repo := NewContentModerationRepo()
		eventService := events.NewEventService(txRunner, events.NewEventRepo(), 5)
		service := NewContentModerationService(tx, repo, eventService, reviewService)

		// Act
		err := service.ReportContent(ctx, "non-existent-review-id", Review, "Inappropriate content", "user-id")

		// Assert
		require.NotNil(t, err)
		require.Equal(t, ErrContentElementDoesntExist, err)
	})
}

func TestReportReview_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		// DI
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), tx, 5)
		repo := NewContentModerationRepo()
		eventService := events.NewEventService(txRunner, events.NewEventRepo(), 5)
		service := NewContentModerationService(tx, repo, eventService, reviewService)

		// Seed data
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)
		ids := testutil.SeedRatingsAndReviews(tx, bookIDs[0], userIDs[0], 4.5, 1.0)

		// Act
		err := service.ReportContent(ctx, ids[1], Review, "Inappropriate content", userIDs[0])

		// Assert
		require.Nil(t, err)

		event, err := eventService.DequeueEvent(ctx)

		require.Nil(t, err)
		require.NotNil(t, event)
		require.Equal(t, EventReviewReported, event.Type)
		require.Equal(t, ids[1], event.AggregateID)
	})
}
