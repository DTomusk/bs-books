package search

import (
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/queries"

	"github.com/gin-gonic/gin"
)

type BookSearchHandler struct {
	service *BookSearchService
}

func NewSearchHandler(service *BookSearchService) *BookSearchHandler {
	return &BookSearchHandler{service: service}
}

// SearchBooks godoc
// @Summary Search Books
// @Description Search for books by title
// @Tags Books
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} BookSearchResponse
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /books/search [get]
func (h *BookSearchHandler) SearchBooks(c *gin.Context) {
	var req BookSearchRequest
	logger := logging.FromContext(c.Request.Context())
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}
	req.PagedRequest.Normalise()

	ctx := c.Request.Context()
	page, err := h.service.SearchBooks(ctx, req.Query, req.PagedRequest.Page, req.PagedRequest.PageSize)

	if err != nil {
		logger.Error("Failed to search books", "error", err)
		c.JSON(500, response.NewInternalServerError("Failed to search books"))
		return
	}
	c.JSON(200, response.Success[[]queries.BookSearchItem]{
		Data: page.Items,
		Meta: response.PageMeta{
			TotalItems: page.Total,
			TotalPages: page.TotalPages,
			Page:       page.Page,
			PageSize:   page.Size,
		},
	})
}
