package main

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/books"
	"bs-books-api/internal/config"
	"bs-books-api/internal/delivery"
	"bs-books-api/internal/queries"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/users"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "bs-books-api/docs"

	_ "github.com/lib/pq"
)

// @title BS Books API
// @version 1.0
// @description This is the API that will change the world
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Enter your JWT token in the format: Bearer <token>
func main() {
	// Load env variables
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Connect to and ping database on startup
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// DI for routes
	userRepo := users.NewUserRepo()
	userService := users.NewUserService(db, userRepo)
	userHandler := users.NewUserHandler(userService)

	jwtService := auth.NewJWTService(cfg.JWT_SECRET_KEY, cfg.JWT_EXPIRATION_MINUTES)
	authService := auth.NewAuthService(db, userService, jwtService)
	authHandler := auth.NewAuthHandler(authService)

	bookReader := queries.NewBookReader(db)
	bookHandler := books.NewBookHandler(bookReader)

	ratingRepo := ratings.NewRatingRepo()
	ratingService := ratings.NewRatingService(db, ratingRepo)
	ratingHandler := ratings.NewRatingHandler(ratingService)

	r := delivery.NewRouter(
		authHandler,
		bookHandler,
		ratingHandler,
		userHandler,
		jwtService,
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		log.Printf("Starting server on port \n 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
