package queries

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetReviewsByBookID_InvalidBookID_Errors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		reader := NewReviewReader(tx)
		ctx := context.Background()

		// Act
		_, err := reader.GetReviewsByBookIDQuery(ctx, "invalid-book-id", 1, 10, 0)

		// Assert
		require.Error(t, err)
	})
}

func TestGetReviewsByBookID_BookWithNoReviews_EmptyPage(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		reader := NewReviewReader(tx)
		ctx := context.Background()
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)

		// Act
		page, err := reader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, page)
		require.Equal(t, 0, page.Total)
		require.Equal(t, 0, len(page.Items))
	})
}

func TestGetReviewsByBookID_BookWithReviews_ReturnsPage(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		reader := NewReviewReader(tx)
		ctx := context.Background()
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		testutil.SeedRatingsAndReviews(tx)

		// Act
		page, err := reader.GetReviewsByBookIDQuery(ctx, bookIDs[0], 1, 10, 0)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, page)
		require.Equal(t, 2, page.Total)
		require.Equal(t, 2, len(page.Items))
		require.Equal(t, "Great book!", page.Items[0].Text)
	})
}
