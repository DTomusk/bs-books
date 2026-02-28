package reviews

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateReview_NoTextThrows(t *testing.T) {
	// Arrange/Act
	review, err := newReview("some-rating-id", "")

	// Assert
	require.Nil(t, review)
	require.Equal(t, ErrEmptyReviewText, err)
}

func TestCreateReview_WithTextSucceeds(t *testing.T) {
	// Arrange/Act
	review, err := newReview("some-rating-id", "This is a review text.")

	// Assert
	require.NotNil(t, review)
	require.NoError(t, err)
	require.Equal(t, "some-rating-id", review.RatingID)
	require.Equal(t, "This is a review text.", review.Text)
	require.NotEmpty(t, review.ID)
}
