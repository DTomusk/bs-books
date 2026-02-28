package books

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateBookWithAuthors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authorIDs := testutil.SeedAuthors(tx)
		repo := NewBooksRepo()
		ctx := context.Background()
		service := NewBooksService(nil, repo)
		book, err := NewBook("Test Book with Authors", authorIDs, "http://example.com/image.jpg", "This is a test synopsis.")
		require.NoError(t, err)

		// Act
		err = service.createBookWithAuthors(book, tx, ctx)

		// Assert
		require.NoError(t, err)

		// Verify book was created correctly
		repoBook, err := repo.getBookByID(book.ID, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, book.ID, repoBook.ID)
		require.Equal(t, book.Title, repoBook.Title)
		require.ElementsMatch(t, book.AuthorIDs, repoBook.AuthorIDs)
	})
}
