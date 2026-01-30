package ratings

type RatingRequest struct {
	BookID     string  `json:"book_id"`
	UserID     string  `json:"user_id"`
	HeartScore float64 `json:"heart_score"`
	PooScore   float64 `json:"poo_score"`
	Review     string  `json:"review,omitempty"`
}
