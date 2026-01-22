package books

import (
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
