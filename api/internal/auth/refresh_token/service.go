package refresh_token

import (
	"bs-books-api/internal/db"
	"context"
	"crypto/rand"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	db              db.DBTX
	tokenExpiryDays int
	hasher          *TokenHasher
	repo            *RefreshTokenRepo
}

func NewRefreshTokenService(db db.DBTX, tokenExpiryDays int, hasher *TokenHasher, repo *RefreshTokenRepo) *RefreshTokenService {
	return &RefreshTokenService{
		db:              db,
		tokenExpiryDays: tokenExpiryDays,
		hasher:          hasher,
		repo:            repo,
	}
}

// Creates a new refresh token in a new family for a new session
// Called at login
func (s *RefreshTokenService) NewSession(ctx context.Context, userID string, ipAddress string) (*RefreshToken, error) {
	// Create a new refresh token
	refreshToken, err := s.createNewToken(userID, ipAddress)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveRefreshToken(ctx, s.db, refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) createNewToken(userID, ipAddress string) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	token := string(b)
	tokenHash := s.hasher.Hash(token)
	return &RefreshToken{
		ID:        uuid.NewString(),
		Token:     token,
		TokenHash: tokenHash,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(time.Duration(s.tokenExpiryDays) * 24 * time.Hour).Unix(),
		UserID:    userID,
		IPAddress: ipAddress,
		FamilyID:  uuid.NewString(),
	}, nil
}
