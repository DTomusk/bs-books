package reviews

import "bs-books-api/internal/delivery/request"

type ReviewListingRequest struct {
	BookID       string `json:"book_id"`
	PagedRequest request.PagedRequest
}
