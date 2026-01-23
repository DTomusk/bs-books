package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DB_URL                 string
	JWT_SECRET_KEY         string
	JWT_EXPIRATION_MINUTES int
	ENV                    string
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

	return &Config{
		DB_URL:                 dbURL,
		JWT_SECRET_KEY:         jwtSecretKey,
		JWT_EXPIRATION_MINUTES: jwtExpirationMinutes,
		ENV:                    env,
	}, nil
}
