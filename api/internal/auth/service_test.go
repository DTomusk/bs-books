package auth

import (
	"bs-books-api/internal/auth/refresh_token"
	"bs-books-api/internal/testutil"
	"bs-books-api/internal/users"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, userService, jwtService, refreshTokenService)
		ctx := context.Background()

		// Act
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")

		// Assert
		require.NoError(t, err)
	})
}

func TestRegister_WeakPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, nil, nil, refreshTokenService)
		ctx := context.Background()
		err := testService.Register(ctx, "blah", "blah@mail.com", "123")
		require.Error(t, err)
		require.Equal(t, ErrShortPassword, err)
	})
}

func TestRegister_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, nil, nil, refreshTokenService)
		ctx := context.Background()
		err := testService.Register(ctx, "user", "invalid-email", "strongpassword")
		require.Error(t, err)
		require.Equal(t, ErrInvalidEmail, err)
	})
}

func TestRegister_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, userService, nil, refreshTokenService)
		ctx := context.Background()
		err := testService.Register(ctx, "testuser", "test@email.com", "password123")
		require.NoError(t, err)

		// Act
		err = testService.Register(ctx, "anotheruser", "test@email.com", "anotherpassword")

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrEmailAlreadyInUse, err)
	})
}

func TestRegister_DuplicateUsername(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, userService, nil, refreshTokenService)
		ctx := context.Background()
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")
		require.NoError(t, err)

		// Act
		err = testService.Register(ctx, "testuser", "different@email.com", "anotherpassword")

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrUsernameAlreadyInUse, err)
	})
}

func TestLogin_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, userService, jwtService, refreshTokenService)
		ctx := context.Background()
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")
		require.NoError(t, err)

		// Act
		token, refreshToken, err := testService.Login(ctx, "test@example.com", "password123", "127.0.0.1")

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotNil(t, refreshToken)
		require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
		require.False(t, refreshToken.IsRevoked)
		require.NotEmpty(t, refreshToken.FamilyID)
		require.NotEmpty(t, refreshToken.TokenHash)
		require.NotEqual(t, refreshToken.TokenHash, refreshToken.Token)
	})
}

func TestLogin_WrongEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, userService, jwtService, refreshTokenService)
		ctx := context.Background()

		// Act
		_, _, err := testService.Login(ctx, "wrong@example.com", "password123", "127.0.0.1")

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}

func TestLogin_WrongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		txRunner := testutil.NewTestTxRunner(tx)
		refreshTokenService := refresh_token.NewRefreshTokenService(txRunner, 7, refresh_token.NewTokenHasher("abc"), refresh_token.NewRefreshTokenRepo())
		testService := NewAuthService(tx, userService, jwtService, refreshTokenService)
		ctx := context.Background()
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")
		require.NoError(t, err)

		// Act
		_, _, err = testService.Login(ctx, "test@example.com", "wrongpassword", "127.0.0.1")

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}
