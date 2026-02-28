package reviews

import "bs-books-api/internal/delivery/response"

type ReviewListingResponse struct {
	Reviews []ReviewListingItem `json:"data"`
	Meta    response.PageMeta   `json:"meta"`
}

type ReviewListingItem struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Username   string  `json:"username"`
	HeartScore float64 `json:"heart_score"`
	PooScore   float64 `json:"poo_score"`
	Text       string  `json:"text"`
}
