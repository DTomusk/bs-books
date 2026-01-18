package delivery

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/books"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/users"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	authHandler *auth.AuthHandler,
	bookHandler *books.BookHandler,
	ratingHandler *ratings.RatingHandler,
	userHandler *users.UserHandler,
	jwtService *auth.JWTService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
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
		booksRoutes.GET("", bookHandler.GetBooks)
	}

	ratingsRoutes := api.Group("/ratings")
	{
		ratingsRoutes.POST("", ratingHandler.CreateRating)
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
