package internal_test

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/books"
	"bs-books-api/internal/delivery/response"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/testutil"
	"bs-books-api/internal/users"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupCreateReviewRouter(
	jwtService *auth.JWTService,
	authService *auth.AuthService,
	userHandler *users.UserHandler,
	ratingService *ratings.RatingService,
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

	ratingsHandler := ratings.NewRatingHandler(ratingService)
	ratingsRoutes := api.Group("/ratings")
	ratingsRoutes.Use(auth.AuthMiddleware(jwtService))
	{
		ratingsRoutes.POST("", ratingsHandler.CreateRating)
	}

	return r
}

func TestRegisterLoginPostRatingWithReview_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// DI
		txRunner := testutil.NewTestTxRunner(tx)
		jwtService := auth.NewJWTService("test_secret_key", 15)
		userRepo := users.NewUserRepo()
		userService := users.NewUserService(tx, userRepo)
		userHandler := users.NewUserHandler(userService)
		authService := auth.NewAuthService(tx, userService, jwtService)
		ratingRepo := ratings.NewRatingRepo()
		bookRepo := books.NewBooksRepo()
		bookService := books.NewBooksService(txRunner, bookRepo)
		reviewRepo := reviews.NewReviewRepo()
		reviewService := reviews.NewReviewService(reviewRepo)
		ratingService := ratings.NewRatingService(txRunner, ratingRepo, bookService, reviewService, nil)
		router := setupCreateReviewRouter(jwtService, authService, userHandler, ratingService)

		// Seed authors and books
		testutil.SeedAuthors(tx)
		bookIDs := testutil.SeedBooks(tx)

		// 1. Register
		body := []byte(`{
		"email": "peepee@poopoo.com",
		"password": "securepassword"
		}`)

		registerReq := jsonRequest("POST", "/api/auth/register", body)
		registerW := httptest.NewRecorder()

		router.ServeHTTP(registerW, registerReq)

		require.Equal(t, 201, registerW.Code)

		// 2. Login
		loginReq := jsonRequest("POST", "/api/auth/login", body)
		loginW := httptest.NewRecorder()

		router.ServeHTTP(loginW, loginReq)

		require.Equal(t, 200, loginW.Code)

		var loginResp response.Success[string]
		err := json.Unmarshal(loginW.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		require.NotEmpty(t, loginResp.Data)

		token := loginResp.Data

		// 3. Post Rating with Review
		ratingBody := []byte(`{
			"book_id": "` + bookIDs[0] + `",
			"heart_score": 4.5,
			"poo_score": 1.0,
			"review": "I HATE READING!!!"
		}`)

		ratingReq := jsonRequest("POST", "/api/ratings", ratingBody)
		ratingReq.Header.Set("Authorization", "Bearer "+token)
		ratingW := httptest.NewRecorder()

		router.ServeHTTP(ratingW, ratingReq)

		require.Equal(t, 201, ratingW.Code)
	})
}
