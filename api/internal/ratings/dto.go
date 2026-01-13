package ratings

type RatingRequest struct {
	BookID     string `json:"book_id"`
	HeartScore int    `json:"heart_score"`
	PooScore   int    `json:"poo_score"`
}
