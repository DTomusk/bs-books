package reviews

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoCreate(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewReviewRepo()
		ctx := context.Background()
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)
		ratingID := testutil.SeedRating(tx, bookIDs[0], userIDs[0], 4.5, 1.0)
		review, err := newReview(ratingID, "This is a test review.")

		// Act
		err = repo.create(review, ctx, tx)

		// Assert
		require.NoError(t, err)
	})
}
