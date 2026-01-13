package books

type BookResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	AuthorID string `json:"author_id"`
}

func ToBookResponse(b *Book) *BookResponse {
	return &BookResponse{
		ID:       b.ID,
		Title:    b.Title,
		AuthorID: b.AuthorID,
	}
}

func ToBookResponses(books []*Book) []*BookResponse {
	responses := make([]*BookResponse, len(books))
	for i, b := range books {
		responses[i] = ToBookResponse(b)
	}
	return responses
}
