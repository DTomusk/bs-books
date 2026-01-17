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

func (h *UserHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, response.NewError("unauthorized", "unauthorized"))
		return
	}
	user, err := h.service.GetUserByID(userID.(string), ctx)
	if err != nil {
		c.JSON(500, response.NewError("internal_server_error", err.Error()))
		return
	}
	if user == nil {
		c.JSON(404, response.NewError("not_found", "user not found"))
		return
	}
	// TODO: define UserResponse struct
	c.JSON(200, gin.H{"id": user.ID, "email": user.Email})
}
