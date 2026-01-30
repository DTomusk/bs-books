package reviews

type Review struct {
	BookID string
	UserID string
	// Note: each review is associated with a rating
	RatingID string
	Text     string
}

func newReview(bookID, userID, ratingID, text string) (*Review, error) {
	if text == "" {
		return nil, ErrEmptyReviewText
	}
	return &Review{
		BookID:   bookID,
		UserID:   userID,
		RatingID: ratingID,
		Text:     text,
	}, nil
}
