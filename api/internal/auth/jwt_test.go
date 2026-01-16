package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseJWT_RoundTrip(t *testing.T) {
	service := NewJWTService("test-secret", 5)
	userID := "user-123"

	token, err := service.generateJWT(userID)
	require.NoError(t, err)

	claims, err := service.parseJWT(token)
	require.NoError(t, err)

	require.Equal(t, userID, claims.UserID)
	require.Greater(t, claims.ExpiresAt.Unix(), claims.IssuedAt.Unix())
}

func TestParseJWT_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret", 5)
	invalidToken := "this.is.not.a.valid.token"

	_, err := service.parseJWT(invalidToken)
	require.ErrorIs(t, err, ErrInvalidToken)
}
