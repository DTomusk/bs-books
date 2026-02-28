package ratings

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"
)

func TestRepoCreateRating(t *testing.T) {
	// Arrange
	r := NewRatingRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)
		userIDs := testutil.SeedUsers(tx)

		ctx := context.Background()

		// Act
		rating, err := newRating(
			bookIDs[0],
			userIDs[0],
			4.5,
			2.0,
		)

		// Assert
		if err != nil {
			t.Fatalf("Failed to create rating entity: %v", err)
		}

		err = r.create(rating, ctx, tx)
		if err != nil {
			t.Fatalf("Failed to create rating: %v", err)
		}
	})
}
