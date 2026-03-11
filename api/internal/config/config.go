package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DB_URL                      string
	JWT_SECRET_KEY              string
	JWT_EXPIRATION_MINUTES      int
	ENV                         string
	GOOGLE_BOOKS_API_KEY        string
	EVENTS_MAX_RETRIES          int
	EVENTS_RETRY_DELAY_SECONDS  int
	AUTHOR_SIMILARITY_THRESHOLD float64
	REVIEW_VISIBILITY_THRESHOLD int
	REFRESH_TOKEN_EXPIRY_DAYS   int
	REFRESH_TOKEN_HASH_SALT     string
	CORS_ALLOWED_ORIGIN         string
}

var (
	ErrMissingDBURL = fmt.Errorf("missing DB_URL environment variable")
)

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return nil, ErrMissingDBURL
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return nil, fmt.Errorf("missing JWT_SECRET_KEY environment variable")
	}

	jwtExpirationMinutesStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	if jwtExpirationMinutesStr == "" {
		return nil, fmt.Errorf("missing JWT_EXPIRATION_MINUTES environment variable")
	}

	jwtExpirationMinutes, err := strconv.Atoi(jwtExpirationMinutesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_MINUTES value: %v", err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	googleBooksAPIKey := os.Getenv("GOOGLE_BOOKS_API_KEY")

	eventsMaxRetriesStr := os.Getenv("EVENTS_MAX_RETRIES")
	eventsMaxRetries, err := strconv.Atoi(eventsMaxRetriesStr)
	if err != nil || eventsMaxRetriesStr == "" {
		return nil, fmt.Errorf("invalid or missing EVENTS_MAX_RETRIES value: %v", err)
	}

	eventsRetryDelaySecondsStr := os.Getenv("EVENTS_RETRY_DELAY_SECONDS")
	eventsRetryDelaySeconds, err := strconv.Atoi(eventsRetryDelaySecondsStr)
	if err != nil || eventsRetryDelaySecondsStr == "" {
		return nil, fmt.Errorf("invalid or missing EVENTS_RETRY_DELAY_SECONDS value: %v", err)
	}

	authorSimilarityThresholdStr := os.Getenv("AUTHOR_SIMILARITY_THRESHOLD")
	authorSimilarityThreshold, err := strconv.ParseFloat(authorSimilarityThresholdStr, 64)
	if err != nil || authorSimilarityThresholdStr == "" {
		return nil, fmt.Errorf("invalid or missing AUTHOR_SIMILARITY_THRESHOLD value: %v", err)
	}

	reviewVisibilityThresholdStr := os.Getenv("REVIEW_VISIBILITY_THRESHOLD")
	reviewVisibilityThreshold, err := strconv.Atoi(reviewVisibilityThresholdStr)

	if err != nil || reviewVisibilityThresholdStr == "" {
		return nil, fmt.Errorf("invalid or missing REVIEW_VISIBILITY_THRESHOLD value: %v", err)
	}

	refreshTokenExpiryDaysStr := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS")
	refreshTokenExpiryDays, err := strconv.Atoi(refreshTokenExpiryDaysStr)
	if err != nil || refreshTokenExpiryDaysStr == "" {
		return nil, fmt.Errorf("invalid or missing REFRESH_TOKEN_EXPIRY_DAYS value: %v", err)
	}

	refreshTokenHashSalt := os.Getenv("REFRESH_TOKEN_HASH_SALT")
	if refreshTokenHashSalt == "" {
		return nil, fmt.Errorf("missing REFRESH_TOKEN_HASH_SALT environment variable")
	}

	corsAllowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if corsAllowedOrigin == "" {
		corsAllowedOrigin = "http://localhost:9000"
	}

	return &Config{
		DB_URL:                      dbURL,
		JWT_SECRET_KEY:              jwtSecretKey,
		JWT_EXPIRATION_MINUTES:      jwtExpirationMinutes,
		ENV:                         env,
		GOOGLE_BOOKS_API_KEY:        googleBooksAPIKey,
		EVENTS_MAX_RETRIES:          eventsMaxRetries,
		EVENTS_RETRY_DELAY_SECONDS:  eventsRetryDelaySeconds,
		AUTHOR_SIMILARITY_THRESHOLD: authorSimilarityThreshold,
		REVIEW_VISIBILITY_THRESHOLD: reviewVisibilityThreshold,
		REFRESH_TOKEN_EXPIRY_DAYS:   refreshTokenExpiryDays,
		REFRESH_TOKEN_HASH_SALT:     refreshTokenHashSalt,
		CORS_ALLOWED_ORIGIN:         corsAllowedOrigin,
	}, nil
}
