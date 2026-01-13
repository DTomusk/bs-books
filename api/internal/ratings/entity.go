package ratings

type Rating struct {
	ID         string
	BookID     string
	HeartScore float64
	PooScore   float64
}

func newRating(bookID string, heartScore, pooScore float64) (*Rating, error) {
	if heartScore < 0 || pooScore < 0 {
		return nil, ErrNegativeScore
	}

	if heartScore > 5.0 || pooScore > 5.0 {
		return nil, ErrLargeScore
	}

	return &Rating{
		BookID:     bookID,
		HeartScore: heartScore,
		PooScore:   pooScore,
	}, nil
}
