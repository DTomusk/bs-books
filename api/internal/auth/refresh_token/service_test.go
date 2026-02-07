package refresh_token

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	})
}
