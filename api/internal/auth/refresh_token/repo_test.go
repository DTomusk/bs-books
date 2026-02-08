package refresh_token

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSaveRefreshToken_ThrowsIfExpired(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), NewRefreshTokenRepo())
		ctx := context.Background()
		userIDs := testutil.SeedUsers(tx)

		existingToken, err := refreshTokenService.createNewToken(userIDs[0], "127.0.0.1")
		require.NoError(t, err)
		// Manually expire token
		existingToken.ExpiresAt = time.Now().Add(-time.Hour).Unix()

		// Act
		err = refreshTokenService.repo.SaveNewRefreshToken(ctx, txRunner.DB(), existingToken)

		// Assert
		require.Error(t, err)
	})
}

func TestSaveNewToken_RevokeFamily_GetToken_ExpectedProperties(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		repo := NewRefreshTokenRepo()
		ctx := context.Background()
		userIDs := testutil.SeedUsers(tx)
		service := NewRefreshTokenService(txRunner, 7, NewTokenHasher("abc"), repo)

		token, err := service.createNewToken(userIDs[0], "127.0.0.1")
		require.NoError(t, err)

		err = repo.SaveNewRefreshToken(ctx, txRunner.DB(), token)
		require.NoError(t, err)

		// Act
		err = repo.RevokeRefreshTokensForFamily(ctx, tx, token.FamilyID)
		require.NoError(t, err)

		fetchedToken, err := repo.GetRefreshTokenByHash(ctx, txRunner.DB(), token.TokenHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, fetchedToken)
		require.Equal(t, token.TokenHash, fetchedToken.TokenHash)
		require.Equal(t, token.UserID, fetchedToken.UserID)
		require.Equal(t, token.IPAddress, fetchedToken.IPAddress)
		require.Equal(t, token.FamilyID, fetchedToken.FamilyID)
		require.Equal(t, token.ExpiresAt, fetchedToken.ExpiresAt)
		require.True(t, fetchedToken.IsRevoked)
	})
}
