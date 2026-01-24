package books

import (
	"bs-books-api/internal/logging"
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

		// Verify book was created correctly
		repoBook, err := repo.getBookByID(book.ID, ctx, tx)
		require.NoError(t, err)
		require.Equal(t, book.ID, repoBook.ID)
		require.Equal(t, book.Title, repoBook.Title)
		require.ElementsMatch(t, book.AuthorIDs, repoBook.AuthorIDs)
	})
}

func TestCreateBookFromExternal_AllAuthorsPresent(t *testing.T) {
	// Arrange
	externalBook := externalBookModel{
		Title:   "External Book",
		Authors: []string{"Author A", "Author B"},
	}
	authorNameIDs := map[string]string{
		"Author A": "author-a-id",
		"Author B": "author-b-id",
	}
	logger := logging.FromContext(context.Background())

	// Act
	book, err := createBookFromExternal(externalBook, authorNameIDs, logger)

	// Assert
	require.NoError(t, err)
	require.Equal(t, externalBook.Title, book.Title)
	require.ElementsMatch(t, []string{"author-a-id", "author-b-id"}, book.AuthorIDs)
}

func TestCreateBookFromExternal_MissingAuthors(t *testing.T) {
	// Arrange
	externalBook := externalBookModel{
		Title:   "External Book",
		Authors: []string{"Author A", "Author C"},
	}
	authorNameIDs := map[string]string{
		"Author A": "author-a-id",
		"Author B": "author-b-id",
	}
	logger := logging.FromContext(context.Background())

	// Act
	book, err := createBookFromExternal(externalBook, authorNameIDs, logger)

	// Assert
	require.Nil(t, book)
	require.Equal(t, ErrNotAllAuthorsPresent, err)
}
