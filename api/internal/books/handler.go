package books

import (
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/queries"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	reader *queries.BookReader
}

func NewBookHandler(r *queries.BookReader) *BookHandler {
	return &BookHandler{reader: r}
}

// GetBooks godoc
// @Summary Get list of books
// @Description Get a list of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {array} BookListResponse
// @Router /books [get]
func (h *BookHandler) GetBooks(ctx *gin.Context) {
	bookDTOs, err := h.reader.GetAllBooksQuery(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, response.NewError("internal_error", "Failed to retrieve books"))
		return
	}
	ctx.JSON(200, response.Success[[]*queries.BookResponse]{Data: bookDTOs})
}
