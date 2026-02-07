package auth

import (
	"bs-books-api/internal/db"
	"context"
)

type RefreshTokenRepo struct{}

func NewRefreshTokenRepo() *RefreshTokenRepo {
	return &RefreshTokenRepo{}
}

func (r *RefreshTokenRepo) RevokeRefreshTokensForUser(userID string) error {
	return nil
}

// Creates a brand new refresh token that isn't descended from another
// Hence family_id is omitted
func (r *RefreshTokenRepo) CreateNewRefreshToken(
	ctx context.Context,
	db db.DBTX,
	tokenHash string,
	token *RefreshToken,
) error {
	var query = `
		INSERT INTO refresh_tokens 
		(id, user_id, token_hash, expires_at, ip_address)
		VALUES ($1, $2, $3, to_timestamp($4), $5)
	`
	_, err := db.ExecContext(ctx, query, token.ID, token.UserID, tokenHash, token.ExpiresAt, token.IPAddress)
	return err
}
