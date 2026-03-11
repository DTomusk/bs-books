package delivery

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/books"
	"bs-books-api/internal/books/search"
	"bs-books-api/internal/content_moderation"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/users"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	authHandler *auth.AuthHandler,
	bookSearchHandler *search.BookSearchHandler,
	ratingHandler *ratings.RatingHandler,
	userHandler *users.UserHandler,
	reviewHandler *reviews.ReviewHandler,
	jwtService *auth.JWTService,
	bookHandler *books.BookHandler,
	contentModerationHandler *content_moderation.ContentModerationHandler,
	corsAllowedOrigin string,
) *gin.Engine {
	r := gin.New()

	r.Use(logging.RequestLoggerMiddleware)
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsAllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	booksRoutes := api.Group("/books")
	{
		booksRoutes.GET("/:id", bookHandler.GetBookByID)

		searchRoutes := booksRoutes.Group("/search")
		{
			searchRoutes.GET("", bookSearchHandler.SearchBooks)
		}

		reviewRoutes := booksRoutes.Group("/:id/reviews")
		{
			reviewRoutes.GET("", reviewHandler.GetReviewsByBookID)
		}
	}

	contentModerationRoutes := api.Group("/moderation")
	{
		protected := contentModerationRoutes.Group("")
		protected.Use(auth.AuthMiddleware(jwtService))
		{
			protected.POST("/report", contentModerationHandler.ReportContent)
		}
	}

	ratingsRoutes := api.Group("/ratings")
	{
		protected := ratingsRoutes.Group("")
		protected.Use(auth.AuthMiddleware(jwtService))
		{
			protected.POST("", ratingHandler.CreateRating)
		}
	}

	usersRoutes := api.Group("/users")
	{
		protected := usersRoutes.Group("")
		protected.Use(auth.AuthMiddleware(jwtService))
		{
			protected.GET("/me", userHandler.GetMe)
		}
	}

	return r
}
