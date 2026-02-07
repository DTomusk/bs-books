package main

import (
	"bs-books-api/internal/auth"
	"bs-books-api/internal/authors"
	"bs-books-api/internal/books"
	"bs-books-api/internal/books/extraction"
	"bs-books-api/internal/books/search"
	"bs-books-api/internal/config"
	"bs-books-api/internal/content_moderation"
	"bs-books-api/internal/db"
	"bs-books-api/internal/delivery"
	"bs-books-api/internal/events"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/queries"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/reviews"
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
	database, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}

	if err := database.Ping(); err != nil {
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
	txRunner := db.NewDBTxRunner(database)

	eventService := events.NewEventService(txRunner, events.NewEventRepo(), cfg.EVENTS_MAX_RETRIES)

	userService := users.NewUserService(database, users.NewUserRepo())
	userHandler := users.NewUserHandler(userService)

	jwtService := auth.NewJWTService(cfg.JWT_SECRET_KEY, cfg.JWT_EXPIRATION_MINUTES)
	refreshTokenService := auth.NewRefreshTokenService(cfg.REFRESH_TOKEN_EXPIRY_DAYS, auth.NewTokenHasher(cfg.REFRESH_TOKEN_HASH_SALT))
	authService := auth.NewAuthService(database, auth.NewAuthRepo(), userService, jwtService, refreshTokenService)
	authHandler := auth.NewAuthHandler(authService)

	authorService := authors.NewAuthorsService(database, authors.NewAuthorsRepo(), cfg.AUTHOR_SIMILARITY_THRESHOLD)
	bookReader := queries.NewBookReader(database)
	bookService := books.NewBooksService(txRunner, books.NewBooksRepo())
	bookExtractionService := extraction.NewBookExtractionService(database, extraction.NewGoogleBooksProvider(externalBooksHTTPClient, cfg.GOOGLE_BOOKS_API_KEY), authorService)
	bookSearchService := search.NewBookSearchService(database, bookReader, bookService, search.NewBookSearchRepo(), bookExtractionService)
	searchHandler := search.NewSearchHandler(bookSearchService)
	bookHandler := books.NewBookHandler(bookReader)

	reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), database, cfg.REVIEW_VISIBILITY_THRESHOLD)
	reviewReader := queries.NewReviewReader(database)
	reviewHandler := reviews.NewReviewHandler(reviewReader)
	ratingService := ratings.NewRatingService(txRunner, ratings.NewRatingRepo(), bookService, reviewService, eventService)
	ratingHandler := ratings.NewRatingHandler(ratingService)

	contentModerationHandler := content_moderation.NewContentModerationHandler(content_moderation.NewContentModerationService(database, content_moderation.NewContentModerationRepo(), eventService, reviewService, userService))

	r := delivery.NewRouter(
		authHandler,
		searchHandler,
		ratingHandler,
		userHandler,
		reviewHandler,
		jwtService,
		bookHandler,
		contentModerationHandler,
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
