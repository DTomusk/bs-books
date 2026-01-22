package authors

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessExternalAuthor_CreatesNewAuthor(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo)
		ctx := context.Background()
		authorName := "New Author Name"

		// Act
		returnedID, err := service.processExternalAuthor(authorName, ctx)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, returnedID)

		// Verify author created in DB
		storedID, err := repo.getIDByName(authorName, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, returnedID, storedID)
	})
}

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

func TestProcessExternalAuthor_AliasMatchReturnsID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo)
		ctx := context.Background()
		author := NewAuthor("Original Author Name")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)
		err = repo.createAuthorAlias(author.ID, "Alias Author Name", ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("Alias Author Name", ctx)

		// Assert
		require.NoError(t, err)
		require.Equal(t, author.ID, returnedID)
	})
}

func TestProcessExternalAuthor_NormalisedMatchCreatesAlias(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo)
		ctx := context.Background()
		author := NewAuthor("J. G. Ballard")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("J G Ballard", ctx)

		// Assert
		require.NoError(t, err)
		require.Equal(t, author.ID, returnedID)

		// Check alias created
		aliasID, err := repo.getIDByAlias("J G Ballard", ctx, tx)
		require.NoError(t, err)
		require.Equal(t, author.ID, aliasID)
	})
}
