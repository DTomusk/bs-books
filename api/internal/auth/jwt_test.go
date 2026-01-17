package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseJWT_RoundTrip(t *testing.T) {
	service := NewJWTService("test-secret", 5)
	userID := "user-123"

	token, err := service.GenerateJWT(userID)
	require.NoError(t, err)

	claims, err := service.ParseJWT(token)
	require.NoError(t, err)

	require.Equal(t, userID, claims.UserID)
	require.Greater(t, claims.ExpiresAt.Unix(), claims.IssuedAt.Unix())
}

func TestParseJWT_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret", 5)
	invalidToken := "this.is.not.a.valid.token"

	_, err := service.ParseJWT(invalidToken)
	require.ErrorIs(t, err, ErrInvalidToken)
}
