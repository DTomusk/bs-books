package search

import "bs-books-api/internal/delivery/response"

type BookSearchResponse struct {
	Data []BookSearchItem  `json:"data"`
	Meta response.PageMeta `json:"meta"`
}

type BookSearchItem struct {
	ID         string             `json:"id"`
	Title      string             `json:"title"`
	Similarity float64            `json:"similarity"`
	ImageURL   string             `json:"image_url"`
	Authors    []AuthorSearchItem `json:"authors"`
}

type AuthorSearchItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
