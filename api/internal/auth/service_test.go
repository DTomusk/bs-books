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
		testService := NewAuthService(tx, userService)
		ctx := context.Background()

		// Act
		err := testService.Register(ctx, "test@example.com", "password123")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

// Test duplicate email
// Test invalid email format
// Test weak password
func TestRegister_WeakPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		testService := NewAuthService(tx, nil)
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
