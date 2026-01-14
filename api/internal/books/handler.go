package books

import (
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/queries"

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
	// TODO: Pass actual db transaction
	bookDTOs := queries.GetAllBooksQuery(nil)
	ctx.JSON(200, response.Success[[]*queries.BookResponse]{Data: bookDTOs})
}
