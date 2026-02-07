package auth

import (
	"crypto/rand"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID string
	// Raw token used in cookie
	Token     string
	IsRevoked bool
	// ExpiresAt is a unix timestamp in seconds
	ExpiresAt int64
	UserID    string
	IPAddress string
}

func NewRefreshToken(userID string, expiresInDays int, ipAddress string) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		ID:        uuid.NewString(),
		Token:     string(b),
		IsRevoked: false,
		ExpiresAt: time.Now().Add(time.Duration(expiresInDays) * 24 * time.Hour).Unix(),
		UserID:    userID,
		IPAddress: ipAddress,
	}, nil
}

type RefreshTokenService struct {
	tokenExpiryDays int
	hasher          *TokenHasher
}

func NewRefreshTokenService(tokenExpiryDays int, hasher *TokenHasher) *RefreshTokenService {
	return &RefreshTokenService{
		tokenExpiryDays: tokenExpiryDays,
		hasher:          hasher,
	}
}

func (s *RefreshTokenService) CreateAndRotate(userID string, ipAddress string) (*RefreshToken, error) {
	return nil, nil
}
