package delivery

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/books"
	"bs-books-api/internal/books/search"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/reviews"
	"bs-books-api/internal/users"
	"net/http"

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
) *gin.Engine {
	r := gin.New()
	r.Use(logging.RequestLoggerMiddleware)
	r.Use(gin.Recovery())

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
