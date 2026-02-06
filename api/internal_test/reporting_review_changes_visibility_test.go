package internal_test

import (
	"bs-books-api/internal/content_moderation"
	"bs-books-api/internal/events"
	"bs-books-api/internal/queries"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReportReview_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		// DI
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), tx, 1)
		repo := content_moderation.NewContentModerationRepo()
		eventService := events.NewEventService(txRunner, events.NewEventRepo(), 5)
		service := content_moderation.NewContentModerationService(tx, repo, eventService, reviewService)
		reviewReader := queries.NewReviewReader(tx)

		// Seed data
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)
		ids := testutil.SeedRatingsAndReviews(tx, bookIDs[0], userIDs[0], 4.5, 1.0)

		// Get reviews for the book
		reviews, err := reviewReader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)
		require.Nil(t, err)
		require.NotNil(t, reviews)
		require.Equal(t, 1, len(reviews.Items))

		// Act
		err = service.ReportContent(ctx, ids[1], content_moderation.Review, "Inappropriate content", userIDs[0])

		// Assert
		require.Nil(t, err)

		// Dequeue event to check it was raised successfully
		event, err := eventService.DequeueEvent(ctx)
		require.Nil(t, err)
		require.NotNil(t, event)
		require.Equal(t, content_moderation.EventReviewReported, event.Type)
		require.Equal(t, ids[1], event.AggregateID)

		// Process event
		err = reviewService.HandleReviewReported(ctx, tx, ids[1])

		// Get book reviews again to check that review is no longer visible
		reviews, err = reviewReader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)
		require.Nil(t, err)
		require.NotNil(t, reviews)
		require.Equal(t, 0, len(reviews.Items))
	})
}

func TestReportReview_HigherThreshold_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		// DI
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), tx, 2)
		repo := content_moderation.NewContentModerationRepo()
		eventService := events.NewEventService(txRunner, events.NewEventRepo(), 5)
		service := content_moderation.NewContentModerationService(tx, repo, eventService, reviewService)
		reviewReader := queries.NewReviewReader(tx)

		// Seed data
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)
		ids := testutil.SeedRatingsAndReviews(tx, bookIDs[0], userIDs[0], 4.5, 1.0)

		// Get reviews for the book
		reviews, err := reviewReader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)
		require.Nil(t, err)
		require.NotNil(t, reviews)
		require.Equal(t, 1, len(reviews.Items))

		// Act: first content report
		err = service.ReportContent(ctx, ids[1], content_moderation.Review, "Inappropriate content", userIDs[0])

		// Assert
		require.Nil(t, err)

		// Dequeue event to check it was raised successfully
		event, err := eventService.DequeueEvent(ctx)
		require.Nil(t, err)
		require.NotNil(t, event)
		require.Equal(t, content_moderation.EventReviewReported, event.Type)
		require.Equal(t, ids[1], event.AggregateID)

		// Process event
		err = reviewService.HandleReviewReported(ctx, tx, ids[1])

		// Get book reviews again to check that review is still visible
		reviews, err = reviewReader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)
		require.Nil(t, err)
		require.NotNil(t, reviews)
		require.Equal(t, 1, len(reviews.Items))

		// Act 2: another user reports the same review
		err = service.ReportContent(ctx, ids[1], content_moderation.Review, "Inappropriate content", userIDs[1])

		// Assert
		require.Nil(t, err)

		// Dequeue event to check it was raised successfully
		event, err = eventService.DequeueEvent(ctx)
		require.Nil(t, err)
		require.NotNil(t, event)
		require.Equal(t, content_moderation.EventReviewReported, event.Type)
		require.Equal(t, ids[1], event.AggregateID)

		// Process event
		err = reviewService.HandleReviewReported(ctx, tx, ids[1])

		// Get book reviews again to check that review is no longer visible
		reviews, err = reviewReader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)
		require.Nil(t, err)
		require.NotNil(t, reviews)
		require.Equal(t, 0, len(reviews.Items))
	})
}
