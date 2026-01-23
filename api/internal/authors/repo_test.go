package authors

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAuthor_NoDuplicate(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		ctx := context.Background()

		author := NewAuthor("Unique Author Name")

		// Act
		err := repo.createAuthor(author, ctx, tx)

		// Assert
		require.NoError(t, err)

		// Verify author was created correctly
		repoAuthor, err := repo.getAuthorByID(author.ID, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, author.ID, repoAuthor.ID)
		require.Equal(t, author.Name, repoAuthor.Name)
		require.Equal(t, author.NormalisedName, repoAuthor.NormalisedName)
		require.Equal(t, author.DuplicateID, repoAuthor.DuplicateID)
	})
}

func TestCreateAuthor_WithDuplicate(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		ctx := context.Background()
		author := NewAuthor("Mr. Bob Write")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		duplicateID := author.ID
		dupeAuthor := NewAuthorWithDuplicate("Mr. Bob Writes", duplicateID)

		// Act
		err = repo.createAuthor(dupeAuthor, ctx, tx)

		// Assert
		require.NoError(t, err)

		// Verify author was created correctly
		repoAuthor, err := repo.getAuthorByID(dupeAuthor.ID, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, dupeAuthor.ID, repoAuthor.ID)
		require.Equal(t, dupeAuthor.Name, repoAuthor.Name)
		require.Equal(t, dupeAuthor.NormalisedName, repoAuthor.NormalisedName)
		require.Equal(t, dupeAuthor.DuplicateID, repoAuthor.DuplicateID)
	})
}

func TestSearchByNormalisedName_ExactMatch(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		ctx := context.Background()
		author := NewAuthor("Pee pee poo poo")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		foundAuthor, err := repo.searchByNormalisedName(author.NormalisedName, ctx, tx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, foundAuthor)
		require.Equal(t, author.ID, foundAuthor.ID)
		require.Equal(t, author.Name, foundAuthor.Name)
		require.Equal(t, author.NormalisedName, foundAuthor.NormalisedName)
	})
}

func TestSearchByNormalisedName_CloseMatch(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		repo := NewAuthorsRepo()
		ctx := context.Background()
		author := NewAuthor("Dr. Seuss")
		err := repo.createAuthor(author, ctx, tx)
		require.NoError(t, err)

		// Act
		foundAuthor, err := repo.searchByNormalisedName("dr seus", ctx, tx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, foundAuthor)
		require.Equal(t, author.ID, foundAuthor.ID)
		require.Equal(t, author.Name, foundAuthor.Name)
		require.Equal(t, author.NormalisedName, foundAuthor.NormalisedName)
	})
}
