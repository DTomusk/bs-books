package search

import "bs-books-api/internal/delivery/request"

type BookSearchRequest struct {
	Query        string
	PagedRequest request.PagedRequest
}
