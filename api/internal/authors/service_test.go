package authors

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessExternalAuthor_MatchReturnsID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo)
		ctx := context.Background()
		author := NewAuthor("My Favourite Author!")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("My Favourite Author!", ctx)

		// Assert
		require.NoError(t, err)
		require.Equal(t, author.ID, returnedID)
	})
}
