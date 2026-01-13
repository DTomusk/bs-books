package books

import (
	"bs-books-api/internal/delivery/response"

	"github.com/gin-gonic/gin"
)

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
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
	books := []*Book{
		{ID: "1", Title: "Book 1", AuthorID: "1"},
		{ID: "2", Title: "Book 2", AuthorID: "2"},
		{ID: "3", Title: "Book 3", AuthorID: "3"},
	}
	bookDTOs := ToBookResponses(books)
	ctx.JSON(200, response.Success[[]*BookResponse]{Data: bookDTOs})
}
