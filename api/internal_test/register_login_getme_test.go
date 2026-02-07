package internal_test

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/testutil"
	"bs-books-api/internal/users"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupAuthRouter(
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
		// DI
		jwtService := auth.NewJWTService("test_secret_key", 15)
		userRepo := users.NewUserRepo()
		userService := users.NewUserService(tx, userRepo)
		userHandler := users.NewUserHandler(userService)
		refreshTokenService := auth.NewRefreshTokenService(7, auth.NewTokenHasher("abc"))
		authService := auth.NewAuthService(tx, auth.NewAuthRepo(), userService, jwtService, refreshTokenService)
		router := setupAuthRouter(jwtService, authService, userHandler)

		// 1. Register
		registerBody := []byte(`{
		"username": "testuser",
		"email": "test@example.com",
		"password": "securepassword"
		}`)

		registerReq := jsonRequest("POST", "/api/auth/register", registerBody)
		registerW := httptest.NewRecorder()

		router.ServeHTTP(registerW, registerReq)

		require.Equal(t, 201, registerW.Code)

		// 2. Login
		loginBody := []byte(`{
		"email": "test@example.com",	
		"password": "securepassword"
		}`)
		loginReq := jsonRequest("POST", "/api/auth/login", loginBody)
		loginW := httptest.NewRecorder()

		router.ServeHTTP(loginW, loginReq)

		require.Equal(t, 200, loginW.Code)

		var loginResp response.Success[string]
		err := json.Unmarshal(loginW.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		require.NotEmpty(t, loginResp.Data)

		token := loginResp.Data

		// 3. Get Me
		meReq := jsonRequest("GET", "/api/users/me", nil)
		meReq.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, meReq)

		require.Equal(t, 200, w.Code)

		var meResp response.Success[users.UserResponse]
		err = json.Unmarshal(w.Body.Bytes(), &meResp)
		require.NoError(t, err)
		require.Equal(t, "test@example.com", meResp.Data.Email)
	})
}

func TestGetMe_NoAuthHeader(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		jwtService := auth.NewJWTService("test_secret_key", 15)
		refreshTokenService := auth.NewRefreshTokenService(7, auth.NewTokenHasher("abc"))
		router := setupAuthRouter(
			jwtService,
			auth.NewAuthService(tx, auth.NewAuthRepo(), users.NewUserService(tx, users.NewUserRepo()), jwtService, refreshTokenService),
			users.NewUserHandler(users.NewUserService(tx, users.NewUserRepo())),
		)
		meReq := jsonRequest("GET", "/api/users/me", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, meReq)

		require.Equal(t, 401, w.Code)
	})
}

func TestGetMe_InvalidToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		jwtService := auth.NewJWTService("test_secret_key", 15)
		refreshTokenService := auth.NewRefreshTokenService(7, auth.NewTokenHasher("abc"))
		router := setupAuthRouter(
			jwtService,
			auth.NewAuthService(tx, auth.NewAuthRepo(), users.NewUserService(tx, users.NewUserRepo()), jwtService, refreshTokenService),
			users.NewUserHandler(users.NewUserService(tx, users.NewUserRepo())),
		)
		meReq := jsonRequest("GET", "/api/users/me", nil)
		meReq.Header.Set("Authorization", "Bearer invalidtoken")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, meReq)

		require.Equal(t, 401, w.Code)
	})
}

// TODO: could move to testutil/somewhere shared
func jsonRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
