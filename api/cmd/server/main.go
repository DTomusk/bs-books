package main

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/config"
	"bs-books-api/internal/delivery"
	"bs-books-api/internal/ratings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "bs-books-api/docs"
)

// @title BS Books API
// @version 1.0
// @description This is the API that will change the world
// @host localhost:8080
// @BasePath /api
func main() {
	// Load env variables
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// TODO: remove
	fmt.Println("Database URL:", cfg.DB_URL)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	bookHandler := books.NewBookHandler()

	ratingService := ratings.NewRatingService()
	ratingHandler := ratings.NewRatingHandler(ratingService)

	r := delivery.NewRouter(bookHandler, ratingHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port \n 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
