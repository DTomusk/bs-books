package books

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryGoogleAPI(t *testing.T) {
	// Arrange
	httpClient := http.DefaultClient
	provider := NewGoogleBooksProvider(httpClient)
	ctx := context.Background()

	// Act
	volumes, err := provider.queryAPI("harry+potter", 1, ctx)

	// Assert
	require.NoError(t, err)
	require.Greater(t, len(volumes), 0)
	require.NotEmpty(t, volumes[0].VolumeInfo.Title)
	require.Greater(t, len(volumes[0].VolumeInfo.Authors), 0)
}

func TestMapGoogleToExternal(t *testing.T) {
	// Arrange
	googleVolumes := []googleVolume{
		{
			VolumeInfo: googleVolumeInfo{
				Title:   "Test Book",
				Authors: []string{"Author One", "Author Two"},
			},
		},
		{
			VolumeInfo: googleVolumeInfo{
				Title:   "Another Book",
				Authors: []string{"Single Author"},
			},
		},
	}

	// Act
	externalBooks := mapGoogleToExternal(googleVolumes)

	// Assert
	require.Len(t, externalBooks, 2)
	require.Equal(t, "Test Book", externalBooks[0].Title)
	require.Equal(t, "Author One", externalBooks[0].Authors[0])
	require.Equal(t, "Author Two", externalBooks[0].Authors[1])
	require.Equal(t, "Another Book", externalBooks[1].Title)
	require.Equal(t, "Single Author", externalBooks[1].Authors[0])
}

func TestSearchBooks(t *testing.T) {
	// Arrange
	httpClient := http.DefaultClient
	provider := NewGoogleBooksProvider(httpClient)
	ctx := context.Background()

	// Act
	externalBooks, err := provider.SearchBooks("lord+of+the+rings", 1, ctx)

	// Assert
	require.NoError(t, err)
	require.Greater(t, len(externalBooks), 0)
	require.NotEmpty(t, externalBooks[0].Title)
	require.NotEmpty(t, externalBooks[0].Authors)
}
