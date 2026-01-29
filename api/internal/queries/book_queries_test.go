package queries

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchBooks_OneExactResult(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)
		ctx := context.Background()
		reader := NewBookReader(tx)

		// Act
		books, err := reader.SearchBooksQuery("big fists", 1, 1, 0, ctx)

		// Assert
		require.NoError(t, err)
		require.Len(t, books.Items, 1)
		require.Equal(t, "Big Fists", books.Items[0].Title)
		require.Equal(t, 1, books.Total)
		require.Equal(t, 1, books.Page)
		require.Equal(t, 1, books.Size)
		require.Equal(t, 1, books.TotalPages)
		require.Equal(t, 1.0, books.Items[0].Similarity)
	})
}
