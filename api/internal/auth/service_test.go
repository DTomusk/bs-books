package auth

import (
	"bs-books-api/internal/testutil"
	"bs-books-api/internal/users"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()

		// Act
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")

		// Assert
		require.NoError(t, err)
	})
}

func TestRegister_WeakPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewAuthService(tx, nil, nil)
		ctx := context.Background()
		err := testService.Register(ctx, "blah", "blah@mail.com", "123")
		require.Error(t, err)
		require.Equal(t, ErrShortPassword, err)
	})
}

func TestRegister_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewAuthService(tx, nil, nil)
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
		testService := NewAuthService(tx, userService, nil)
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

func TestLogin_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")
		require.NoError(t, err)

		// Act
		token, err := testService.Login(ctx, "test@example.com", "password123")

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})
}

func TestLogin_WrongEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()

		// Act
		_, err := testService.Login(ctx, "wrong@example.com", "password123")

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
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()
		err := testService.Register(ctx, "testuser", "test@example.com", "password123")
		require.NoError(t, err)

		// Act
		_, err = testService.Login(ctx, "test@example.com", "wrongpassword")

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}
