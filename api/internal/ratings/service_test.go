package ratings

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceCreateRating(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRatingRepo()

		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())

		reviewService := reviews.NewReviewService(reviews.NewReviewRepo())

		testService := NewRatingService(txRunner, r, bookService, reviewService)
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)

		ctx := context.Background()

		// Act
		err := testService.CreateRating(
			bookIDs[0],
			userIDs[0],
			4.5,
			2.0,
			"",
			ctx,
		)

		// Assert
		require.NoError(t, err)
	})
}

func TestServiceCreateRating_BookNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRatingRepo()
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo())
		testService := NewRatingService(txRunner, r, bookService, reviewService)
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)
		userIds := testutil.SeedUsers(tx)

		ctx := context.Background()

		// Act
		err := testService.CreateRating(
			"non-existent-book-id",
			userIds[0],
			4.5,
			2.0,
			"",
			ctx,
		)

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrBookNotFound, err)
	})
}

func TestServiceCreateRating_UserNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRatingRepo()
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo())
		testService := NewRatingService(txRunner, r, bookService, reviewService)
		testutil.SeedAuthors(tx)
		bookIds := testutil.SeedBooks(tx)
		testutil.SeedUsers(tx)
		ctx := context.Background()

		// Act
		err := testService.CreateRating(
			bookIds[0],
			"non-existent-user-id",
			4.5,
			2.0,
			"",
			ctx,
		)

		// Assert
		require.Error(t, err)
	})
}

func TestServiceCreateRating_RatingAlreadyExists(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRatingRepo()
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo())
		testService := NewRatingService(txRunner, r, bookService, reviewService)
		testutil.SeedAuthors(tx)
		bookIds := testutil.SeedBooks(tx)
		userIds := testutil.SeedUsers(tx)

		ctx := context.Background()

		// First, create a rating
		err := testService.CreateRating(
			bookIds[0],
			userIds[0],
			4.5,
			2.0,
			"",
			ctx,
		)
		require.NoError(t, err)

		// Act - try to create the same rating again
		err = testService.CreateRating(
			bookIds[0],
			userIds[0],
			3.0,
			1.0,
			"",
			ctx,
		)

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrRatingAlreadyExists, err)
	})
}

func TestServiceCreateRating_SameBookDifferentUsers(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRatingRepo()
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo())
		testService := NewRatingService(txRunner, r, bookService, reviewService)
		testutil.SeedAuthors(tx)
		bookIds := testutil.SeedBooks(tx)
		userIds := testutil.SeedUsers(tx)
		ctx := context.Background()

		// First user creates a rating
		err := testService.CreateRating(
			bookIds[0],
			userIds[0],
			4.5,
			2.0,
			"",
			ctx,
		)
		require.NoError(t, err)

		// Act - second user creates a rating for the same book
		err = testService.CreateRating(
			bookIds[0],
			userIds[1],
			3.5,
			1.0,
			"",
			ctx,
		)

		// Assert
		require.NoError(t, err)
	})
}

func TestServiceCreateRating_DifferentBooksSameUser(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRatingRepo()
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		reviewService := reviews.NewReviewService(reviews.NewReviewRepo())
		testService := NewRatingService(txRunner, r, bookService, reviewService)
		testutil.SeedAuthors(tx)
		bookIds := testutil.SeedBooks(tx)
		userIds := testutil.SeedUsers(tx)
		ctx := context.Background()

		// First rating for the first book
		err := testService.CreateRating(
			bookIds[0],
			userIds[0],
			4.5,
			2.0,
			"",
			ctx,
		)
		require.NoError(t, err)

		// Act - create a rating for a different book by the same user
		err = testService.CreateRating(
			bookIds[1],
			userIds[0],
			3.5,
			1.0,
			"",
			ctx,
		)

		// Assert
		require.NoError(t, err)
	})
}
