package books

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBook(t *testing.T) {
	// Arrange
	title := "Test Book"
	authorIDs := []string{"author1", "author2"}
	synopsis := "This is a test synopsis."
	imageURL := "http://example.com/image.jpg"

	// Act
	book, err := NewBook(title, authorIDs, imageURL, synopsis)

	// Assert
	require.NoError(t, err)
	require.Equal(t, title, book.Title)
	require.Equal(t, authorIDs, book.AuthorIDs)
	require.Equal(t, synopsis, book.Synopsis)
	require.Equal(t, imageURL, book.ImageURL)
}

func TestNewBook_NoAuthors(t *testing.T) {
	// Arrange
	title := "Test Book"
	var authorIDs []string

	// Act
	book, err := NewBook(title, authorIDs, "", "")

	// Assert
	require.Nil(t, book)
	require.EqualError(t, err, "no authors provided for the book")
}

func TestNewBook_EmptyAuthors(t *testing.T) {
	// Arrange
	title := "Test Book"
	authorIDs := []string{}

	// Act
	book, err := NewBook(title, authorIDs, "", "")

	// Assert
	require.Nil(t, book)
	require.EqualError(t, err, "no authors provided for the book")
}
