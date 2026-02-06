package main

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/config"
	"bs-books-api/internal/content_moderation"
	"bs-books-api/internal/db"
	"bs-books-api/internal/events"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/ratings"
	"bs-books-api/internal/reviews"
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
	eventService := events.NewEventService(txRunner, eventRepo, cfg.EVENTS_MAX_RETRIES)

	bookService := books.NewBooksService(txRunner, books.NewBooksRepo())

	reviewService := reviews.NewReviewService(reviews.NewReviewRepo(), database, cfg.REVIEW_VISIBILITY_THRESHOLD)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Starting event processor...")
		for {
			select {
			case <-time.After(time.Duration(cfg.EVENTS_RETRY_DELAY_SECONDS) * time.Second):
			case <-ctx.Done():
				slog.Info("Event processor shutting down")
				return
			default:
			}

			event, err := eventService.DequeueEvent(ctx)
			if err != nil {
				slog.Error("Failed to dequeue event", "error", err)
				continue
			}

			if event == nil {
				slog.Info("No event to process")
				time.Sleep(time.Duration(cfg.EVENTS_RETRY_DELAY_SECONDS) * time.Second)
				continue
			}

			switch event.Type {
			case ratings.EventTypeRatingCreated:
				err = txRunner.WithTx(ctx, func(tx *sql.Tx) error {
					bookID := event.AggregateID

					var payload ratings.RatingCreatedPayload
					err := json.Unmarshal(event.Payload, &payload)
					if err != nil {
						slog.Error("Failed to unmarshal event payload", "error", err, "eventID", event.ID)
						return err
					}

					heartScore := payload.HeartScore
					pooScore := payload.PooScore

					slog.Info("Processing rating created event", "bookID", bookID, "heartScore", heartScore, "pooScore", pooScore)
					err = bookService.AddRatingToBook(ctx, tx, bookID, heartScore, pooScore)
					if err != nil {
						slog.Error("Failed to add rating to book", "error", err, "bookID", bookID)
						return err
					}

					// Mark event as processed
					err = eventService.MarkEventProcessed(ctx, tx, event.ID)
					if err != nil {
						slog.Error("Failed to mark event as processed", "error", err, "eventID", event.ID)
						return err
					}

					slog.Info("Successfully processed rating created event", "eventID", event.ID)

					return nil
				})
				if err != nil {
					slog.Error("Failed to process rating created event", "error", err, "eventID", event.ID)
					continue
				}
			case content_moderation.EventReviewReported:
				err = txRunner.WithTx(ctx, func(tx *sql.Tx) error {
					reviewID := event.AggregateID

					slog.Info("Processing review reported event", "reviewID", reviewID)
					// Call review service to increment reports on the specified review and change the visiblity if relevant
					err := reviewService.HandleReviewReported(ctx, tx, reviewID)
					if err != nil {
						slog.Error("Failed to handle review reported event", "error", err, "reviewID", reviewID)
						return err
					}

					err = eventService.MarkEventProcessed(ctx, tx, event.ID)
					if err != nil {
						slog.Error("Failed to mark event as processed", "error", err, "eventID", event.ID)
						return err
					}

					slog.Info("Successfully processed review reported event", "eventID", event.ID)

					return nil
				})
				if err != nil {
					slog.Error("Failed to process review reported event", "error", err, "eventID", event.ID)
					continue
				}
			case content_moderation.EventUserReported:
				err = txRunner.WithTx(ctx, func(tx *sql.Tx) error {
					userID := event.AggregateID

					slog.Info("Processing user reported event", "userID", userID)

					// TODO: implement handling, for now just mark as processed
					err = eventService.MarkEventProcessed(ctx, tx, event.ID)
					if err != nil {
						slog.Error("Failed to mark event as processed", "error", err, "eventID", event.ID)
						return err
					}

					slog.Info("Successfully processed user reported event", "eventID", event.ID)

					return nil
				})
			default:
				slog.Warn("Unknown event type", "type", event.Type, "eventID", event.ID)
			}
		}
	}()

	<-quitCh
	slog.Info("Shutting down...")
}
