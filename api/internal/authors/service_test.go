package authors

import (
	"bs-books-api/internal/logging"
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
		service := NewAuthorsService(tx, repo, 0.8)
		ctx := context.Background()
		authorName := "New Author Name"

		// Act
		returnedID, err := service.processExternalAuthor(authorName, ctx, logging.FromContext(ctx))

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
		service := NewAuthorsService(tx, repo, 0.8)
		ctx := context.Background()
		author := NewAuthor("My Favourite Author!")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("My Favourite Author!", ctx, logging.FromContext(ctx))

		// Assert
		require.NoError(t, err)
		require.Equal(t, author.ID, returnedID)
	})
}

func TestProcessExternalAuthor_AliasMatchReturnsID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo, 0.8)
		ctx := context.Background()
		author := NewAuthor("Original Author Name")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)
		err = repo.createAuthorAlias(author.ID, "Alias Author Name", ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("Alias Author Name", ctx, logging.FromContext(ctx))

		// Assert
		require.NoError(t, err)
		require.Equal(t, author.ID, returnedID)
	})
}

func TestProcessExternalAuthor_NormalisedMatchCreatesAlias(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo, 0.8)
		ctx := context.Background()
		author := NewAuthor("J. G. Ballard")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("J G Ballard", ctx, logging.FromContext(ctx))

		// Assert
		require.NoError(t, err)
		require.Equal(t, author.ID, returnedID)

		// Check alias created
		aliasID, err := repo.getIDByAlias("J G Ballard", ctx, tx)
		require.NoError(t, err)
		require.Equal(t, author.ID, aliasID)
	})
}

func TestProcessExternalAuthor_SimilarNormalisedNameCreatesPotentialDuplicate(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo, 0.7)
		ctx := context.Background()
		author := NewAuthor("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		returnedID, err := service.processExternalAuthor("ABCDEFGHIJKLMNOPQRSTUVWXY", ctx, logging.FromContext(ctx))

		// Assert
		require.NoError(t, err)
		require.NotEqual(t, author.ID, returnedID)

		// Verify new author created with duplicate ID set
		newAuthor, err := repo.getAuthorByID(returnedID, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, "ABCDEFGHIJKLMNOPQRSTUVWXY", newAuthor.Name)
		require.Equal(t, author.ID, *newAuthor.DuplicateID)
	})
}

func TestProcessExternalAuthors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		service := NewAuthorsService(tx, repo, 0.7)
		ctx := context.Background()
		authorNames := []string{
			"J R R Tolkien",
			"J. R. R. Tolkien",
			"Charles Bukowski",
			"Charles bukowsky",
			"Pee Pee Poo Poo",
		}

		// Act
		authorIDs := service.ProcessExternalAuthors(authorNames, ctx)

		// Assert
		require.Len(t, authorIDs, len(authorNames))

		// Verify authors created correctly
		idSet := make(map[string]struct{})
		for _, name := range authorNames {
			id, exists := authorIDs[name]
			require.True(t, exists)
			require.NotEmpty(t, id)
			idSet[id] = struct{}{}
		}
		require.Len(t, idSet, 4, "Expected 4 unique author IDs from 5 names due to Tolkien name variants")

		require.Equal(t, authorIDs["J R R Tolkien"], authorIDs["J. R. R. Tolkien"], "Expected Tolkien name variants to map to same ID")
		require.NotEqual(t, authorIDs["Charles Bukowski"], authorIDs["Charles bukowsky"], "Expected Bukowski name variants to map to different IDs")

		// Verify Tolkien alias
		aliasID, err := repo.getIDByAlias("J. R. R. Tolkien", ctx, tx)
		require.NoError(t, err)
		require.Equal(t, authorIDs["J R R Tolkien"], aliasID)

		// Verify Bukowski duplicate
		dupeAuthor, err := repo.getAuthorByID(authorIDs["Charles bukowsky"], ctx, tx)
		require.NoError(t, err)
		require.NotNil(t, dupeAuthor.DuplicateID)
		originalID, exists := authorIDs["Charles Bukowski"]
		require.True(t, exists)
		require.Equal(t, originalID, *dupeAuthor.DuplicateID)
	})
}
