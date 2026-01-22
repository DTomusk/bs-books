package queries

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllBooksQuery(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		testutil.SeedAuthors(tx)
		testutil.SeedBooks(tx)
		reader := NewBookReader(tx)

		ctx := context.Background()

		// Act
		books, err := reader.GetAllBooksQuery(ctx)

		// Assert
		require.NoError(t, err)
		require.Len(t, books, 2)
	})
}
