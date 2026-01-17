package auth

import (
	"bs-books-api/internal/delivery/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(s *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, response.NewError("missing_auth_header", ErrMissingAuthHeader.Error()))
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(401, response.NewError("invalid_auth_header_format", ErrInvalidAuthHeader.Error()))
			return
		}

		tokenString := authHeader[7:]

		claims, err := s.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, response.NewError("invalid_token", ErrInvalidToken.Error()))
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
