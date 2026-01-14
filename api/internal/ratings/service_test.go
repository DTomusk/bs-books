package ratings

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"
)

func TestServiceCreateRating(t *testing.T) {
	r := NewRatingRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewRatingService(tx, r)
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)

		ctx := context.Background()

		rating, err := testService.CreateRating(
			"23681e21-08d4-43e1-b0b6-8d6f75a9b8b3",
			4.5,
			2.0,
			ctx,
		)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rating == nil {
			t.Fatal("expected rating, got nil")
		}
	})
}

func TestServiceCreateRating_BookNotFound(t *testing.T) {
	r := NewRatingRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewRatingService(tx, r)
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)

		ctx := context.Background()

		_, err := testService.CreateRating(
			"non-existent-book-id",
			4.5,
			2.0,
			ctx,
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		// TODO: Check for specific error type
	})
}
