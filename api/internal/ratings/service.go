package ratings

type RatingService struct{}

func NewRatingService() *RatingService {
	return &RatingService{}
}

func (s *RatingService) CreateRating(bookID string, heartScore int, pooScore int) (*Rating, error) {
	rating := &Rating{
		BookID:     bookID,
		HeartScore: heartScore,
		PooScore:   pooScore,
	}
	return rating, nil
}
