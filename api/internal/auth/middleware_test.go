package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/protected", middleware, func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(200, gin.H{"userID": userID})
	})

	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("test-secret", 5)
	middleware := AuthMiddleware(jwtService)
	router := setupTestRouter(middleware)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, 401, w.Code)
	require.Contains(t, w.Body.String(), ErrMissingAuthHeader.Error())
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("test-secret", 5)
	middleware := AuthMiddleware(jwtService)
	router := setupTestRouter(middleware)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormatToken")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, 401, w.Code)
	require.Contains(t, w.Body.String(), ErrInvalidAuthHeader.Error())
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("test-secret", 5)
	middleware := AuthMiddleware(jwtService)
	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, 401, w.Code)
	require.Contains(t, w.Body.String(), ErrInvalidToken.Error())
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("test-secret", 5)
	token, err := jwtService.generateJWT("user-123")
	require.NoError(t, err)

	router := setupTestRouter(AuthMiddleware(jwtService))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "user-123")
}
