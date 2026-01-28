package books

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateBook(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authorIDs := testutil.SeedAuthors(tx)
		repo := NewBooksRepo()
		ctx := context.Background()
		book, err := NewBook("A New Book", authorIDs, "http://example.com/image.jpg", "This is a synopsis.")
		require.NoError(t, err)

		// Act
		err = repo.createBook(book, ctx, tx)

		// Assert
		require.NoError(t, err)

		// Verify book was created correctly
		repoBook, err := repo.getBookByID(book.ID, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, book.ID, repoBook.ID)
		require.Equal(t, book.Title, repoBook.Title)

		// Insert author associations and verify
		err = repo.addAuthorsToBook(book.ID, book.AuthorIDs, ctx, tx)
		require.NoError(t, err)

		repoBook, err = repo.getBookByID(book.ID, ctx, tx)
		require.NoError(t, err)
		require.ElementsMatch(t, book.AuthorIDs, repoBook.AuthorIDs)
	})
}
