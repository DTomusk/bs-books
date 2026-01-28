package main

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/authors"
	"bs-books-api/internal/books"
	"bs-books-api/internal/books/extraction"
	"bs-books-api/internal/books/search"
	"bs-books-api/internal/config"
	"bs-books-api/internal/delivery"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/queries"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/users"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
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

	logger := logging.New(cfg.ENV)
	slog.SetDefault(logger)

	slog.Info("Configuration loaded", "env", cfg.ENV)

	// Connect to and ping database on startup
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		return
	}

	slog.Info("Connected to database")

	externalBooksHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     60 * time.Second,
		},
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// DI for routes
	userService := users.NewUserService(db, users.NewUserRepo())
	userHandler := users.NewUserHandler(userService)

	jwtService := auth.NewJWTService(cfg.JWT_SECRET_KEY, cfg.JWT_EXPIRATION_MINUTES)
	authService := auth.NewAuthService(db, userService, jwtService)
	authHandler := auth.NewAuthHandler(authService)

	authorService := authors.NewAuthorsService(db, authors.NewAuthorsRepo(), 0.8)

	bookReader := queries.NewBookReader(db)
	bookService := books.NewBooksService(db, books.NewBooksRepo())
	bookExtractionService := extraction.NewBookExtractionService(db, extraction.NewGoogleBooksProvider(externalBooksHTTPClient), authorService)
	bookSearchService := search.NewBookSearchService(db, bookReader, bookService, search.NewBookSearchRepo(), bookExtractionService)
	searchHandler := search.NewSearchHandler(bookSearchService)

	ratingService := ratings.NewRatingService(db, ratings.NewRatingRepo())
	ratingHandler := ratings.NewRatingHandler(ratingService)

	r := delivery.NewRouter(
		authHandler,
		searchHandler,
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
		slog.Info("Starting server on port", "port", 8080)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
