package books

type BookListResponse struct {
	Data []*BookResponse `json:"data"`
}
