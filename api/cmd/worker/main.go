package main

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/config"
	"bs-books-api/internal/db"
	"bs-books-api/internal/events"
	"bs-books-api/internal/logging"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("Connected to database")

	txRunner := db.NewDBTxRunner(database)

	eventRepo := events.NewEventRepo()
	eventService := events.NewEventService(txRunner, eventRepo, 5)

	bookService := books.NewBooksService(txRunner, books.NewBooksRepo())

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Starting event processor...")
		for {
			event, err := eventService.DequeueEvent(ctx)
			if err != nil {
				slog.Error("Failed to dequeue event", "error", err)
				continue
			}

			if event == nil {
				slog.Info("No event to process")
				time.Sleep(5 * time.Second)
				continue
			}

			switch event.Type {
			case "rating.created":
				bookID := event.AggregateID

				// Unmarshal the JSON payload into a map
				var data map[string]float64
				err := json.Unmarshal(event.Payload, &data)
				if err != nil {
					slog.Error("Failed to unmarshal event payload", "error", err, "eventID", event.ID)
					continue
				}

				heartScore, ok := data["heart_score"]
				if !ok {
					slog.Error("Invalid heart_score in event data", "eventID", event.ID)
					continue
				}
				pooScore, ok := data["poo_score"]
				if !ok {
					slog.Error("Invalid poo_score in event data", "eventID", event.ID)
					continue
				}

				slog.Info("Processing rating created event", "bookID", bookID, "heartScore", heartScore, "pooScore", pooScore)
				err = bookService.AddRatingToBook(bookID, heartScore, pooScore, ctx)
				if err != nil {
					slog.Error("Failed to add rating to book", "error", err, "bookID", bookID)
					continue
				}
			}
		}
	}()

	<-quitCh
	slog.Info("Shutting down...")
}
