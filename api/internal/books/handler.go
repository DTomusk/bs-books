package books

import (
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/queries"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	reader *queries.BookReader
}

func NewBookHandler(reader *queries.BookReader) *BookHandler {
	return &BookHandler{reader: reader}
}

// GetBookByID godoc
// @Summary Get Book by ID
// @Description Retrieve book details by its ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} BookDetailsResponse "Successful response"
// @Failure 404 {object} response.ErrorResponse "Book not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	book, err := h.reader.GetBookByID(ctx, id)

	if err != nil {
		c.JSON(500, response.NewInternalServerError("Failed to get book"))
		return
	}

	if book == nil {
		c.JSON(404, response.NewNotFoundError("Book not found"))
		return
	}

	c.JSON(200, response.Success[*queries.BookDetails]{
		Data: book,
	})
}
