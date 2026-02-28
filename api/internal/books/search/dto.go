package search

import "bs-books-api/internal/delivery/request"

type BookSearchRequest struct {
	Query        string `form:"query" binding:"required"`
	PagedRequest request.PagedRequest
}
