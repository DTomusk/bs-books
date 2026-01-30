package ratings

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceCreateRating(t *testing.T) {
	// Arrange
	r := NewRatingRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		testService := NewRatingService(tx, r, bookService)
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)

		ctx := context.Background()

		// Act
		err := testService.CreateRating(
			"23681e21-08d4-43e1-b0b6-8d6f75a9b8b3",
			4.5,
			2.0,
			ctx,
		)

		// Assert
		require.NoError(t, err)
	})
}

func TestServiceCreateRating_BookNotFound(t *testing.T) {
	// Arrange
	r := NewRatingRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
		testService := NewRatingService(tx, r, bookService)
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)

		ctx := context.Background()

		// Act
		err := testService.CreateRating(
			"non-existent-book-id",
			4.5,
			2.0,
			ctx,
		)

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrBookNotFound, err)
	})
}
