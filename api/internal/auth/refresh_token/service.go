package refresh_token

import (
	"bs-books-api/internal/db"
	"context"
	"crypto/rand"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	txRunner        db.TxRunner
	tokenExpiryDays int
	hasher          *TokenHasher
	repo            *RefreshTokenRepo
}

func NewRefreshTokenService(txRunner db.TxRunner, tokenExpiryDays int, hasher *TokenHasher, repo *RefreshTokenRepo) *RefreshTokenService {
	return &RefreshTokenService{
		txRunner:        txRunner,
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

	err = s.repo.SaveRefreshToken(ctx, s.txRunner.DB(), refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) RefreshSession(ctx context.Context, oldRefreshToken string, ipAddress string) (*RefreshToken, error) {
	// Hash token
	tokenHash := s.hasher.Hash(oldRefreshToken)

	// Get token from DB by hash
	oldToken, err := s.repo.GetRefreshTokenByHash(ctx, s.txRunner.DB(), tokenHash)
	if err != nil {
		return nil, err
	}
	if oldToken == nil {
		return nil, ErrInvalidRefreshToken
	}

	var newToken *RefreshToken

	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		// If token is expired or revoked, revoke entire family and return error
		if oldToken.IsRevoked || oldToken.ExpiresAt < time.Now().Unix() {
			err := s.repo.RevokeRefreshTokensForFamily(ctx, tx, oldToken.FamilyID)
			if err != nil {
				return err
			}
			return ErrInvalidRefreshToken
		}

		// Token valid, create new token in same family with current as parent
		newToken, err := s.createChildToken(ctx, oldToken, ipAddress)
		if err != nil {
			return err
		}

		err = s.repo.SetReplacedBy(ctx, tx, oldToken.ID, newToken.ID)
		if err != nil {
			return err
		}

		err = s.repo.SaveRefreshToken(ctx, tx, newToken)
		if err != nil {
			return err
		}
		return nil
	})

	return newToken, nil
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

func (s *RefreshTokenService) createChildToken(ctx context.Context, parent *RefreshToken, ipAddress string) (*RefreshToken, error) {
	childToken, err := s.createNewToken(parent.UserID, ipAddress)
	if err != nil {
		return nil, err
	}
	childToken.FamilyID = parent.FamilyID
	return childToken, nil
}
