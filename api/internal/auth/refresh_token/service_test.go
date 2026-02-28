package refresh_token

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that NewSession creates the expected refresh token
func TestNewSession(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()
		userIDs := testutil.SeedUsers(tx)

		// Act
		token, err := refreshTokenService.NewSession(ctx, userIDs[0], "127.0.0.1")

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, token.Token)
		require.NotEmpty(t, token.TokenHash)
		require.Equal(t, userIDs[0], token.UserID)
		require.NotEqual(t, token.Token, token.TokenHash)
		require.False(t, token.IsRevoked)
		require.Nil(t, token.ReplacedByID)
	})
}

// Assert that the token created by NewSession is persisted correctly and can be fetched by hash
func TestNewSession_FetchPersistedToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()
		userIDs := testutil.SeedUsers(tx)

		// Act
		token, err := refreshTokenService.NewSession(ctx, userIDs[0], "127.0.0.1")
		require.NoError(t, err)

		fetchedToken, err := refreshTokenService.repo.GetRefreshTokenByHash(ctx, txRunner.DB(), token.TokenHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, fetchedToken)
		require.Equal(t, token.TokenHash, fetchedToken.TokenHash)
		require.Equal(t, token.UserID, fetchedToken.UserID)
		require.Equal(t, token.IPAddress, fetchedToken.IPAddress)
		require.Equal(t, token.FamilyID, fetchedToken.FamilyID)
		require.Equal(t, token.ExpiresAt, fetchedToken.ExpiresAt)
	})
}

func TestRefreshSession_NoExistingToken_ReturnsError(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()

		// Act
		newToken, err := refreshTokenService.RefreshSession(ctx, "nonexistenttoken", "127.0.0.1")

		// Assert
		require.ErrorIs(t, err, ErrInvalidRefreshToken)
		require.Nil(t, newToken)
	})
}

func TestRefreshSession_SucceedsAndRevokesOldToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// We start by making an initial session, then we refresh using the resulting token
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()
		userIDs := testutil.SeedUsers(tx)
		oldToken, err := refreshTokenService.NewSession(ctx, userIDs[0], "127.0.0.1")
		require.NoError(t, err)

		// Act
		newToken, err := refreshTokenService.RefreshSession(ctx, oldToken.Token, "127.0.0.1")
		require.NoError(t, err)
		require.NotNil(t, newToken)
		require.NotEqual(t, oldToken.Token, newToken.Token)

		// Assert that the old token is revoked
		fetchedOldToken, err := refreshTokenService.repo.GetRefreshTokenByHash(ctx, txRunner.DB(), oldToken.TokenHash)
		require.NoError(t, err)
		require.NotNil(t, fetchedOldToken)
		require.True(t, fetchedOldToken.IsRevoked)

		// Assert that the new token is valid and has correct properties
		fetchedNewToken, err := refreshTokenService.repo.GetRefreshTokenByHash(ctx, txRunner.DB(), newToken.TokenHash)
		require.NoError(t, err)
		require.NotNil(t, fetchedNewToken)
		require.Equal(t, newToken.TokenHash, fetchedNewToken.TokenHash)
		require.Equal(t, newToken.UserID, fetchedNewToken.UserID)
		require.Equal(t, newToken.IPAddress, fetchedNewToken.IPAddress)
		require.Equal(t, newToken.FamilyID, fetchedNewToken.FamilyID)
		require.Equal(t, newToken.ExpiresAt, fetchedNewToken.ExpiresAt)
		require.False(t, fetchedNewToken.IsRevoked)

		// Assert that they're in the same family and old token is marked as parent of new token
		require.Equal(t, fetchedOldToken.FamilyID, fetchedNewToken.FamilyID)
		require.Equal(t, *fetchedOldToken.ReplacedByID, fetchedNewToken.ID)
	})
}

func TestRevokeSession_Succeeds(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()
		userIDs := testutil.SeedUsers(tx)
		token, err := refreshTokenService.NewSession(ctx, userIDs[0], "127.0.0.1")
		require.NoError(t, err)

		// Act
		err = refreshTokenService.RevokeSession(ctx, token.Token)
		require.NoError(t, err)

		// Assert that the token is revoked
		fetchedToken, err := refreshTokenService.repo.GetRefreshTokenByHash(ctx, txRunner.DB(), token.TokenHash)
		require.NoError(t, err)
		require.NotNil(t, fetchedToken)
		require.True(t, fetchedToken.IsRevoked)
	})
}

func TestRevokeSession_NonexistentToken_ReturnsError(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()

		// Act
		err := refreshTokenService.RevokeSession(ctx, "nonexistenttoken")

		// Assert
		require.ErrorIs(t, err, ErrInvalidRefreshToken)
	})
}
