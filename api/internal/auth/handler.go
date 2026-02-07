package auth

import (
	"bs-books-api/internal/auth/refresh_token"
	"bs-books-api/internal/delivery/response"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param register body AuthRegisterRequest true "Register Request"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req AuthRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}

	err := h.service.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case ErrInvalidEmail, ErrShortPassword, ErrLongPassword:
			c.JSON(400, response.NewError("invalid_request", err.Error()))
			return
		case ErrEmailAlreadyInUse:
			c.JSON(409, response.NewError("email_in_use", err.Error()))
			return
		}
		c.JSON(500, response.NewError("internal_error", err.Error()))
		return
	}

	// TODO: login on successful registration
	c.JSON(201, response.Ok())
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body AuthLoginRequest true "Login Request"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}

	jwt, refreshToken, err := h.service.Login(ctx, req.Email, req.Password, c.ClientIP())
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			c.JSON(401, response.NewError("invalid_credentials", err.Error()))
			return
		}
		c.JSON(500, response.NewInternalServerError(err.Error()))
		return
	}

	// Set http only cookie with refresh token
	c.SetCookie(
		"refresh_token",
		refreshToken.Token,
		int(refreshToken.ExpiresAt-time.Now().Unix()),
		"/",
		"",
		// TODO: set to true in production
		false,
		true,
	)

	c.JSON(200, response.Success[string]{Data: jwt})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh JWT token using refresh token cookie
// @Tags auth
// @Accept json
// @Produce json
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(401, response.NewError("missing_refresh_token", "Refresh token cookie is required"))
		return
	}

	jwt, newRefreshToken, err := h.service.RefreshToken(ctx, refreshToken, c.ClientIP())
	if err != nil {
		switch err {
		case refresh_token.ErrInvalidRefreshToken:
			c.JSON(401, response.NewError("invalid_refresh_token", "Invalid refresh token"))
			return
		}
		c.JSON(500, response.NewInternalServerError(err.Error()))
		return
	}

	c.SetCookie(
		"refresh_token",
		newRefreshToken.Token,
		int(newRefreshToken.ExpiresAt-time.Now().Unix()),
		"/",
		"",
		// TODO: set to true in production
		false,
		true,
	)

	c.JSON(200, response.Success[string]{Data: jwt})
}
