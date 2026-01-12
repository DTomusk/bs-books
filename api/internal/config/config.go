package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB_URL string
}

var (
	ErrMissingDBURL = fmt.Errorf("missing DB_URL environment variable")
)

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return nil, ErrMissingDBURL
	}

	return &Config{
		DB_URL: dbURL,
	}, nil
}
