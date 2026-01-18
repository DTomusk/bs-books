package auth

import (
	"bs-books-api/internal/testutil"
	"bs-books-api/internal/users"
	"context"
	"database/sql"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()

		// Act
		err := testService.Register(ctx, "test@example.com", "password123")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestRegister_WeakPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewAuthService(tx, nil, nil)
		ctx := context.Background()
		err := testService.Register(ctx, "blah@mail.com", "123")
		if err == nil {
			t.Fatal("expected error for weak password, got nil")
		}
		if err != ErrShortPassword {
			t.Fatalf("expected ErrShortPassword, got %v", err)
		}
	})
}

func TestRegister_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewAuthService(tx, nil, nil)
		ctx := context.Background()
		err := testService.Register(ctx, "invalid-email", "strongpassword")
		if err == nil {
			t.Fatal("expected error for invalid email, got nil")
		}
		if err != ErrInvalidEmail {
			t.Fatalf("expected ErrInvalidEmail, got %v", err)
		}
	})
}

func TestRegister_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		testService := NewAuthService(tx, userService, nil)
		ctx := context.Background()
		err := testService.Register(ctx, "test@email.com", "password123")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		// Act
		err = testService.Register(ctx, "test@email.com", "anotherpassword")

		// Assert
		if err == nil {
			t.Fatal("expected error for duplicate email, got nil")
		}
		if err != ErrEmailAlreadyInUse {
			t.Fatalf("expected ErrEmailAlreadyInUse, got %v", err)
		}
	})
}

func TestLogin_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()
		err := testService.Register(ctx, "test@example.com", "password123")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		// Act
		token, err := testService.Login(ctx, "test@example.com", "password123")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Fatal("expected token, got empty string")
		}
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
		if err == nil {
			t.Fatal("expected error for wrong email, got nil")
		}
		if err != ErrInvalidCredentials {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestLogin_WrongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := users.NewUserService(tx, users.NewUserRepo())
		jwtService := NewJWTService("test_secret_key", 15)
		testService := NewAuthService(tx, userService, jwtService)
		ctx := context.Background()
		err := testService.Register(ctx, "test@example.com", "password123")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		// Act
		_, err = testService.Login(ctx, "test@example.com", "wrongpassword")

		// Assert
		if err == nil {
			t.Fatal("expected error for wrong password, got nil")
		}
		if err != ErrInvalidCredentials {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}
