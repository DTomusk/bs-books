package ratings

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"
)

func TestCreateRecipe(t *testing.T) {
	r := NewRatingRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		rating, err := newRating(
			"43681e21-08d4-43e1-b0b6-8d6f75a9b8b1",
			4.5,
			2.0,
		)

		if err != nil {
			t.Fatalf("Failed to create rating entity: %v", err)
		}

		err = r.Create(rating, ctx, tx)
		if err != nil {
			t.Fatalf("Failed to create rating: %v", err)
		}
	})
}
