package search

type BookSearchResponse struct {
	Data []BookSearchItem `json:"data"`
	Meta PageMeta         `json:"meta"`
}

type BookSearchItem struct {
	ID      string             `json:"id"`
	Title   string             `json:"title"`
	Authors []AuthorSearchItem `json:"authors"`
}

type AuthorSearchItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PageMeta struct {
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	Size       int `json:"size"`
}
