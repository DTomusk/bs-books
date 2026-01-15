package users

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"
)

func TestGetUserByEmail_UserExists(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userRepo := NewUserRepo()
		userService := NewUserService(tx, userRepo)
		ctx := context.Background()

		userService.CreateUser("e@mail.com", "beep boop boop", ctx)

		// Act
		retrieved_user, err := userService.GetUserByEmail("e@mail.com", ctx)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if retrieved_user == nil {
			t.Fatal("expected user, got nil")
		}
		if retrieved_user.Email != "e@mail.com" {
			t.Fatalf("expected email to be 'e@mail.com', got %s", retrieved_user.Email)
		}
	})
}
