package refresh_token

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
func (r *RefreshTokenRepo) SaveRefreshToken(
	ctx context.Context,
	db db.DBTX,
	token *RefreshToken,
) error {
	var query = `
		INSERT INTO refresh_tokens 
		(id, user_id, token_hash, expires_at, ip_address, family_id)
		VALUES ($1, $2, $3, to_timestamp($4), $5, $6)
	`
	_, err := db.ExecContext(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.IPAddress, token.FamilyID)
	return err
}
