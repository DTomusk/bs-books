package internal_test

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/events"
	"bs-books-api/internal/queries"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// This test checks that when a rating is created, an event is raised,
// and that event can be processed to update the book's average ratings
func TestCreateRatingRaisesEvent_EventUpdatesBook(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange: DI
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), tx)
		eventService := events.NewEventService(txRunner, events.NewEventRepo(), 5)
		ratingService := ratings.NewRatingService(txRunner, ratings.NewRatingRepo(), bookService, reviewService, eventService)

		bookReader := queries.NewBookReader(tx)

		// Arrange: seed data
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)
		ctx := context.Background()

		// Act 1: create rating
		err := ratingService.CreateRating(
			bookIDs[0],
			userIDs[0],
			4.5,
			2.0,
			"Great book!",
			ctx,
		)

		// Assert rating creation
		require.NoError(t, err)

		// Assert correct event raised
		event, err := eventService.DequeueEvent(ctx)
		require.NoError(t, err)
		require.NotNil(t, event)
		require.Equal(t, ratings.EventTypeRatingCreated, event.Type)
		require.Equal(t, event.AggregateID, bookIDs[0])
		var payload ratings.RatingCreatedPayload
		err = json.Unmarshal(event.Payload, &payload)
		require.NoError(t, err)
		require.Equal(t, 4.5, payload.HeartScore)
		require.Equal(t, 2.0, payload.PooScore)

		// Act 2: process event to update book
		err = bookService.AddRatingToBook(ctx, tx, event.AggregateID, payload.HeartScore, payload.PooScore)
		require.NoError(t, err)

		// Assert book's average ratings updated
		book, err := bookReader.GetBookByID(ctx, bookIDs[0])
		require.NoError(t, err)
		require.Equal(t, 4.5, book.HeartScore)
		require.Equal(t, 2.0, book.PooScore)
		require.Equal(t, 1, book.NumberOfRatings)
	})
}
