package books

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateDedupeKey(t *testing.T) {
	tests := []struct {
		name     string
		book     externalBookModel
		expected string
	}{
		{
			name: "Simple Case",
			book: externalBookModel{
				Title:   " The Great Gatsby ",
				Authors: []string{"F. Scott Fitzgerald"},
			},
			expected: "great gatsby|f scott fitzgerald,",
		},
		{
			name: "Different Punctuation",
			book: externalBookModel{
				Title:   "Moby-Dick!",
				Authors: []string{"  Herman Melville "},
			},
			expected: "moby dick|herman melville,",
		},
		{
			name: "Book series entry",
			book: externalBookModel{
				Title:   "Twilight: New Moon",
				Authors: []string{"Stephenie Meyer"},
			},
			expected: "twilight new moon|stephenie meyer,",
		},
		{
			name: "Multiple authors",
			book: externalBookModel{
				Title:   "Good Omens",
				Authors: []string{"Neil Gaiman", "Terry Pratchett"},
			},
			expected: "good omens|neil gaiman,terry pratchett,",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateDedupeKey(tt.book)
			if got != tt.expected {
				t.Errorf("GenerateDedupeKey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDedupeExternalBooks_OneBook(t *testing.T) {
	// Arrange
	books := []externalBookModel{
		{Title: "Book One", Authors: []string{"Author A"}},
	}

	// Act
	result := deduplicateExternalBooks(books)

	// Assert
	require.Equal(t, 1, len(result), "Expected 1 book after deduplication")
	require.Equal(t, "Book One", result[0].Title, "Book title should match")
	require.Equal(t, 1, len(result[0].Authors), "Book should have 1 author")
	require.Equal(t, "Author A", result[0].Authors[0], "Author name should match")
}

func TestDedupeExternalBooks_MultipleDuplicates(t *testing.T) {
	// Arrange
	books := []externalBookModel{
		{Title: "Book One", Authors: []string{"Author A"}},
		{Title: " Book: One ", Authors: []string{" Author A "}},
		{Title: "Book Two", Authors: []string{"Author B"}},
		{Title: "Book One!", Authors: []string{"Author A"}},
		{Title: "Book Three", Authors: []string{"Author C"}},
	}

	// Act
	result := deduplicateExternalBooks(books)

	// Assert
	require.Equal(t, 3, len(result), "Expected 3 unique books after deduplication")

	expectedTitles := map[string]bool{
		"Book One":   false,
		"Book Two":   false,
		"Book Three": false,
	}

	for _, book := range result {
		if _, exists := expectedTitles[book.Title]; exists {
			expectedTitles[book.Title] = true
		} else {
			t.Errorf("Unexpected book title found: %s", book.Title)
		}
	}

	for title, found := range expectedTitles {
		require.True(t, found, "Expected book title not found: %s", title)
	}
}

func TestDedupeExternalBooks_SameTitlesDifferentAuthors(t *testing.T) {
	// Arrange
	books := []externalBookModel{
		{Title: "Common Title", Authors: []string{"Author A"}},
		{Title: " Common Title ", Authors: []string{" Author B "}},
		{Title: "Common Title!", Authors: []string{"Author A"}},
	}

	// Act
	result := deduplicateExternalBooks(books)

	// Assert
	require.Equal(t, 2, len(result), "Expected 2 unique books after deduplication")
}

func TestDedupeExternalBooks_SameAuthorsDifferentTitles(t *testing.T) {
	// Arrange
	books := []externalBookModel{
		{Title: "Title One", Authors: []string{"Author A"}},
		{Title: " Title Two ", Authors: []string{" Author A "}},
		{Title: "Title One!", Authors: []string{"Author A"}},
	}

	// Act
	result := deduplicateExternalBooks(books)

	// Assert
	require.Equal(t, 2, len(result), "Expected 2 unique books after deduplication")
}

func TestDedupeExternalBooks_EmptyList(t *testing.T) {
	// Arrange
	books := []externalBookModel{}

	// Act
	result := deduplicateExternalBooks(books)

	// Assert
	require.Equal(t, 0, len(result), "Expected 0 books after deduplication of empty list")
}
