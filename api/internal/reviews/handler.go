package reviews

import (
	"bs-books-api/internal/delivery/request"
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/logging"
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
// @Param id path string true "Book ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} ReviewListingResponse
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /books/{id}/reviews [get]
func (h *ReviewHandler) GetReviewsByBookID(c *gin.Context) {
	var uriReq ReviewListingUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}

	var queryReq request.PagedRequest
	if err := c.ShouldBindQuery(&queryReq); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}
	queryReq.Normalise()

	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	logger.Info("Fetching reviews for book", "book_id", uriReq.BookID, "page", queryReq.Page, "page_size", queryReq.PageSize)
	page, err := h.reader.GetReviewsByBookIDQuery(ctx, uriReq.BookID, queryReq.Page, queryReq.PageSize, queryReq.Offset())
	if err != nil {
		logger.Error("Failed to get reviews for book ID", "book_id", uriReq.BookID, "error", err)
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
