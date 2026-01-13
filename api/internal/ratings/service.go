package ratings

type RatingService struct{}

func NewRatingService() *RatingService {
	return &RatingService{}
}

func (s *RatingService) CreateRating(bookID string, heartScore float64, pooScore float64) (*Rating, error) {
	rating, err := newRating(bookID, heartScore, pooScore)
	return rating, err
}
