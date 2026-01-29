package extraction

import (
	"bs-books-api/internal/logging"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateBookFromExternal(t *testing.T) {
	// Arrange
	externalBook := externalBookModel{
		Title:    "Test Book",
		Authors:  []string{"Author One", "Author Two"},
		ImageURL: "http://example.com/image.jpg",
		Synopsis: "This is a test synopsis.",
	}
	authorNameIDs := map[string]string{
		"Author One": "author1",
		"Author Two": "author2",
	}
	logger := logging.FromContext(context.Background())

	// Act
	book, err := createBookFromExternal(externalBook, authorNameIDs, logger)

	// Assert
	require.NoError(t, err)
	require.Equal(t, externalBook.Title, book.Title)
	require.Equal(t, []string{"author1", "author2"}, book.AuthorIDs)
	require.Equal(t, externalBook.ImageURL, book.ImageURL)
	require.Equal(t, externalBook.Synopsis, book.Synopsis)
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
