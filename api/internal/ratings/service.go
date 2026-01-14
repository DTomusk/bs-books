package ratings

type RatingService struct {
	repo *ratingRepo
}

func NewRatingService(r *ratingRepo) *RatingService {
	return &RatingService{
		repo: r,
	}
}

func (s *RatingService) CreateRating(bookID string, heartScore float64, pooScore float64) (*Rating, error) {
	rating, err := newRating(bookID, heartScore, pooScore)

	if err != nil {
		return nil, err
	}

	// TODO: Pass actual context and db transaction
	err = s.repo.create(rating, nil, nil)

	return rating, err
}
