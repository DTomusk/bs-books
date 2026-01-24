package search

import (
	"bs-books-api/internal/delivery/response"

	"github.com/gin-gonic/gin"
)

type BookSearchHandler struct {
	service *BookSearchService
}

func NewSearchHandler(service *BookSearchService) *BookSearchHandler {
	return &BookSearchHandler{service: service}
}

func (h *BookSearchHandler) SearchBooks(c *gin.Context) {
	var req BookSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}
	req.PagedRequest.Normalise()

	ctx := c.Request.Context()
	page, err := h.service.SearchBooks(ctx, req.Query, req.PagedRequest.Page, req.PagedRequest.PageSize)

	if err != nil {
		c.JSON(500, response.NewInternalServerError("Failed to search books"))
		return
	}
	c.JSON(200, response.Success[[]BookResult]{
		Data: page.Items,
		Meta: response.PageMeta{
			TotalItems: page.TotalItems,
			TotalPages: page.TotalPages,
			Page:       page.Page,
			PageSize:   page.Size,
		},
	})
}
