package users

import (
	"bs-books-api/internal/delivery/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetMe godoc
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags users
// @Accept  json
// @Produce  json
// @Router /users/me [get]
// @Security BearerAuth
// @Success 200 {object} GetMeResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
func (h *UserHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, response.ErrUnauthorized)
		return
	}
	user, err := h.service.GetUserByID(userID.(string), ctx)
	if err != nil {
		c.JSON(500, response.NewInternalServerError(err.Error()))
		return
	}
	// If we get this far and user is nil, something went wrong
	if user == nil {
		c.JSON(500, response.NewInternalServerError("user not found"))
		return
	}
	c.JSON(200, response.Success[UserResponse]{Data: UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}})
}
