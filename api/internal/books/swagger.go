package books

type BookDetailsResponse struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	ImageURL string             `json:"image_url"`
	Synopsis string             `json:"synopsis"`
	Authors  []AuthorSearchItem `json:"authors"`
}

type AuthorSearchItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
