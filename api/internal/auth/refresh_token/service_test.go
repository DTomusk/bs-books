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
