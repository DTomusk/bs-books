package books

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBook(t *testing.T) {
	// Arrange
	title := "Test Book"
	authorIDs := []string{"author1", "author2"}

	// Act
	book := NewBook(title, authorIDs)

	// Assert
	require.Equal(t, title, book.Title)
	require.Equal(t, authorIDs, book.AuthorIDs)
}
