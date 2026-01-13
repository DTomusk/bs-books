package delivery

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/ratings"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	bookHandler *books.BookHandler,
	ratingHandler *ratings.RatingHandler,
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

	books := api.Group("/books")
	{
		books.GET("", bookHandler.GetBooks)
	}

	ratings := api.Group("/ratings")
	{
		ratings.POST("", ratingHandler.CreateRating)
	}

	return r
}
