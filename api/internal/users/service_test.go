package users

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail_UserExists(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userRepo := NewUserRepo()
		userService := NewUserService(tx, userRepo)
		ctx := context.Background()

		userService.CreateUser("username", "e@mail.com", "beep boop boop", ctx)

		// Act
		retrieved_user, err := userService.GetUserByEmail("e@mail.com", ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrieved_user)
		require.Equal(t, "e@mail.com", retrieved_user.Email)
		require.Equal(t, "username", retrieved_user.Username)
	})
}
