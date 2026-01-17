package internal_test

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/testutil"
	"bs-books-api/internal/users"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupUsersMeRouter(
	jwtService *auth.JWTService,
	authService *auth.AuthService,
	userHandler *users.UserHandler,
) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api")

	usersRoutes := api.Group("/users")
	protected := usersRoutes.Group("")
	protected.Use(auth.AuthMiddleware(jwtService))
	protected.GET("/me", userHandler.GetMe)

	authHandler := auth.NewAuthHandler(authService)
	authRoutes := api.Group("/auth")
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/register", authHandler.Register)

	return r
}

func TestRegisterLoginMe_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		jwtService := auth.NewJWTService("test_secret_key", 15)
		userRepo := users.NewUserRepo()
		userService := users.NewUserService(tx, userRepo)
		userHandler := users.NewUserHandler(userService)
		authService := auth.NewAuthService(tx, userService, jwtService)
		router := setupUsersMeRouter(jwtService, authService, userHandler)

		// 1. Register
		body := []byte(`{
		"email": "test@example.com",
		"password": "securepassword"
		}`)

		registerReq := httptest.NewRequest(
			"POST",
			"/api/auth/register",
			bytes.NewBuffer(body),
		)
		registerReq.Header.Set("Content-Type", "application/json")
		registerW := httptest.NewRecorder()
		router.ServeHTTP(registerW, registerReq)
		require.Equal(t, 201, registerW.Code)

		// 2. Login
		loginReq := httptest.NewRequest(
			"POST",
			"/api/auth/login",
			bytes.NewBuffer(body),
		)
		loginReq.Header.Set("Content-Type", "application/json")
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)
		require.Equal(t, 200, loginW.Code)

		var loginResp response.Success[string]
		err := json.Unmarshal(loginW.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		require.NotEmpty(t, loginResp.Data)

		token := loginResp.Data

		// 3. Get Me
		req := httptest.NewRequest("GET", "/api/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		fmt.Println(w.Body.String())
		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "test@example.com")
	})
}
