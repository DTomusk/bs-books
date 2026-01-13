package ratings

import (
	"bs-books-api/internal/delivery/response"

	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	service *RatingService
}

func NewRatingHandler(service *RatingService) *RatingHandler {
	return &RatingHandler{service: service}
}

// CreateRating godoc
// @Summary Create a new rating
// @Description Create a new rating for a book
// @Tags ratings
// @Accept  json
// @Produce  json
// @Param rating body RatingRequest true "Rating to create"
// @Success 201
// @Router /ratings [post]
func (h *RatingHandler) CreateRating(ctx *gin.Context) {
	var req RatingRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, response.ErrInvalidRequest)
		return
	}

	_, err := h.service.CreateRating(req.BookID, req.HeartScore, req.PooScore)

	if err != nil {
		ctx.JSON(500, response.NewError("internal_error", "Failed to create rating"))
		return
	}

	ctx.JSON(201, response.Ok())
}
