package authors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAuthorEntity_NoDuplicate(t *testing.T) {
	// Arrange / Act
	author := NewAuthor("Unique Author Name")

	// Assert
	require.NotEmpty(t, author.ID)
	require.Equal(t, "Unique Author Name", author.Name)
	require.Equal(t, "unique author name", author.NormalisedName)
	require.Nil(t, author.DuplicateID)
}

func TestCreateAuthorEntity_WithDuplicate(t *testing.T) {
	// Arrange

	duplicateID := "123e4567-e89b-12d3-a456-426614174000"
	// Act
	author := NewAuthorWithDuplicate("Another Author", duplicateID)

	// Assert
	require.NotEmpty(t, author.ID)
	require.Equal(t, "Another Author", author.Name)
	require.Equal(t, "another author", author.NormalisedName)
	require.NotNil(t, author.DuplicateID)
	require.Equal(t, duplicateID, *author.DuplicateID)
}
