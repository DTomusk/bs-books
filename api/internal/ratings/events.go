package ratings

const (
	EventTypeRatingCreated = "rating.created"
)

type RatingCreatedPayload struct {
	HeartScore float64 `json:"heart_score"`
	PooScore   float64 `json:"poo_score"`
}
