package refresh_token

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
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
	logger := logging.FromContext(ctx)
	logger.Info("Creating new refresh token for user", "userID", userID, "ipAddress", ipAddress)
	refreshToken, err := s.createNewToken(userID, ipAddress)
	if err != nil {
		logger.Error("Failed to create new refresh token", "error", err)
		return nil, err
	}

	err = s.repo.SaveNewRefreshToken(ctx, s.txRunner.DB(), refreshToken)
	if err != nil {
		logger.Error("Failed to save new refresh token", "error", err)
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) RefreshSession(ctx context.Context, oldRefreshToken string, ipAddress string) (*RefreshToken, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Refreshing session with old refresh token", "oldRefreshToken", oldRefreshToken, "ipAddress", ipAddress)

	// Hash token
	tokenHash := s.hasher.Hash(oldRefreshToken)

	// Get token from DB by hash
	oldToken, err := s.repo.GetRefreshTokenByHash(ctx, s.txRunner.DB(), tokenHash)
	if err != nil {
		logger.Error("Failed to fetch refresh token by hash", "error", err)
		return nil, err
	}
	if oldToken == nil {
		logger.Warn("No refresh token found for hash", "tokenHash", tokenHash)
		return nil, ErrInvalidRefreshToken
	}

	var newToken *RefreshToken

	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		// If token is expired or revoked, revoke entire family and return error
		if oldToken.IsRevoked || oldToken.ExpiresAt < time.Now().Unix() {
			logger.Warn("Refresh token is expired or revoked, revoking entire family", "tokenID", oldToken.ID, "familyID", oldToken.FamilyID)
			err := s.repo.RevokeRefreshTokensForFamily(ctx, tx, oldToken.FamilyID)
			if err != nil {
				logger.Error("Failed to revoke refresh token family", "familyID", oldToken.FamilyID, "error", err)
				return err
			}
			logger.Info("Revoked refresh token family due to expired/revoked token usage", "familyID", oldToken.FamilyID)
			return ErrInvalidRefreshToken
		}

		err = s.repo.RevokeRefreshTokensForFamily(ctx, tx, oldToken.FamilyID)
		if err != nil {
			logger.Error("Failed to revoke refresh token family", "familyID", oldToken.FamilyID, "error", err)
			return err
		}
		logger.Info("Revoked refresh token family due to refresh", "familyID", oldToken.FamilyID)

		// Token valid, create new token in same family with current as parent
		newToken, err = s.createChildToken(oldToken, ipAddress)
		if err != nil {
			logger.Error("Failed to create child refresh token", "error", err)
			return err
		}

		err = s.repo.SaveRefreshToken(ctx, tx, newToken)
		if err != nil {
			logger.Error("Failed to save new refresh token", "error", err)
			return err
		}

		err = s.repo.SetReplacedBy(ctx, tx, oldToken.ID, newToken.ID)
		if err != nil {
			logger.Error("Failed to set replaced_by for old refresh token", "oldTokenID", oldToken.ID, "newTokenID", newToken.ID, "error", err)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

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

func (s *RefreshTokenService) createChildToken(parent *RefreshToken, ipAddress string) (*RefreshToken, error) {
	childToken, err := s.createNewToken(parent.UserID, ipAddress)
	if err != nil {
		return nil, err
	}
	childToken.FamilyID = parent.FamilyID
	return childToken, nil
}
