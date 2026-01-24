package books

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractUniqueAuthors(t *testing.T) {
	// Arrange
	books := []externalBookModel{
		{Title: "Book One", Authors: []string{"Author A", "Author B"}},
		{Title: "Book Two", Authors: []string{"Author B", "Author C"}},
		{Title: "Book Three", Authors: []string{"Author A"}},
	}

	// Act
	uniqueAuthors := extractUniqueAuthors(books)

	// Assert
	require.Len(t, uniqueAuthors, 3)
	require.Contains(t, uniqueAuthors, "Author A")
	require.Contains(t, uniqueAuthors, "Author B")
	require.Contains(t, uniqueAuthors, "Author C")
}

func TestCreateBookWithAuthors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authorIDs := testutil.SeedAuthors(tx)
		repo := NewBooksRepo()
		ctx := context.Background()
		service := NewBooksService(nil, repo, nil, nil)
		book, err := NewBook("Test Book with Authors", authorIDs)
		require.NoError(t, err)

		// Act
		err = service.CreateBookWithAuthors(book, tx, ctx)

		// Assert
		require.NoError(t, err)
	})
}
