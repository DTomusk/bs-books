package reviews

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID string
	// Note: each review is associated with a rating
	RatingID  string
	Text      string
	CreatedAt time.Time
}

func newReview(ratingID, text string) (*Review, error) {
	if text == "" {
		return nil, ErrEmptyReviewText
	}
	if len(text) > 500 {
		return nil, ErrReviewTextTooLong
	}
	return &Review{
		ID:        uuid.NewString(),
		RatingID:  ratingID,
		Text:      text,
		CreatedAt: time.Now().UTC(),
	}, nil
}
