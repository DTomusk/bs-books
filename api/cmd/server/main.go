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
	"os"
	"os/signal"
	"syscall"

	_ "bs-books-api/docs"
)

// @title BS Books API
// @version 1.0
// @description This is the API that will change the world
// @host localhost:8080
// @BasePath /api
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	fmt.Println("Database URL:", cfg.DB_URL)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh
	log.Println("Shutting down server...")

	cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
