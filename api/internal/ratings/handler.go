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
func (h *RatingHandler) CreateRating(c *gin.Context) {
	ctx := c.Request.Context()
	var req RatingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}

	err := h.service.CreateRating(req.BookID, req.HeartScore, req.PooScore, ctx)

	switch err {
	case ErrNegativeScore, ErrLargeScore:
		c.JSON(400, response.NewError("invalid_score", err.Error()))
		return
	}

	if err != nil {
		c.JSON(500, response.NewInternalServerError("Failed to create rating"))
		return
	}

	c.JSON(201, response.Ok())
}
