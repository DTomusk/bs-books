package users

import "github.com/gin-gonic/gin"

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	// Implementation for getting the current user's information
	// Get userid from context and pass it to service
}
