package ratings

type RatingRequest struct {
	BookID     string  `json:"book_id"`
	HeartScore float64 `json:"heart_score"`
	PooScore   float64 `json:"poo_score"`
}
