package reviews

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleReviewReported_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		repo := NewReviewRepo()
		service := NewReviewService(repo, tx, 5)

		// Seed data
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)
		ids := testutil.SeedRatingsAndReviews(tx, bookIDs[0], userIDs[0], 4.5, 1.0)

		// Act
		err := service.HandleReviewReported(ctx, tx, ids[1])

		// Assert
		require.Nil(t, err)
	})
}
