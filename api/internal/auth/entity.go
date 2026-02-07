package auth

import (
	"crypto/rand"
	"time"
)

type RefreshToken struct {
	Token     string
	IsRevoked bool
	// ExpiresAt is a unix timestamp in seconds
	ExpiresAt int64
	UserID    string
}

func NewRefreshToken(userID string, expiresInDays int) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		Token:     string(b),
		IsRevoked: false,
		ExpiresAt: time.Now().Add(time.Duration(expiresInDays) * 24 * time.Hour).Unix(),
		UserID:    userID,
	}, nil
}
