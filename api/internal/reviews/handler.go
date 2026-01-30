package reviews

import (
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/queries"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reader *queries.ReviewReader
}

func NewReviewHandler(reader *queries.ReviewReader) *ReviewHandler {
	return &ReviewHandler{reader: reader}
}

// GetReviewsByBookID godoc
// @Summary Get reviews for a book
// @Description Retrieve all reviews associated with a specific book
// @Tags reviews
// @Accept  json
// @Produce  json
// @Param book_id path string true "Book ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} ReviewListingResponse
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /reviews [get]
func (h *ReviewHandler) GetReviewsByBookID(c *gin.Context) {
	var req ReviewListingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}
	req.PagedRequest.Normalise()

	ctx := c.Request.Context()
	page, err := h.reader.GetReviewsByBookIDQuery(ctx, req.BookID, req.PagedRequest.Page, req.PagedRequest.PageSize, req.PagedRequest.Offset())
	if err != nil {
		c.JSON(500, response.NewInternalServerError("Failed to get reviews"))
		return
	}
	c.JSON(200, response.Success[[]queries.ReviewItem]{
		Data: page.Items,
		Meta: response.PageMeta{
			TotalItems: page.Total,
			TotalPages: page.TotalPages,
			Page:       page.Page,
			PageSize:   page.Size,
		},
	})

}
