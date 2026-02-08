package refresh_token

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
	"time"
)

type RefreshTokenRepo struct{}

func NewRefreshTokenRepo() *RefreshTokenRepo {
	return &RefreshTokenRepo{}
}

type RefreshTokenRow struct {
	ID           string
	UserID       string
	TokenHash    string
	ExpiresAt    time.Time
	IPAddress    string
	FamilyID     string
	RevokedAt    sql.NullTime
	ReplacedByID sql.NullString
}

func (r *RefreshTokenRepo) RevokeRefreshTokensForUser(userID string) error {
	return nil
}

// Creates a brand new refresh token that isn't descended from another
// Hence family_id is omitted
func (r *RefreshTokenRepo) SaveNewRefreshToken(
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

func (r *RefreshTokenRepo) GetRefreshTokenByHash(
	ctx context.Context,
	db db.DBTX,
	tokenHash string,
) (*RefreshToken, error) {
	var query = `
		SELECT id, user_id, token_hash, expires_at, ip_address, family_id, revoked_at, replaced_by_token_id
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	row := db.QueryRowContext(ctx, query, tokenHash)
	var tokenRow RefreshTokenRow
	err := row.Scan(&tokenRow.ID, &tokenRow.UserID, &tokenRow.TokenHash, &tokenRow.ExpiresAt, &tokenRow.IPAddress, &tokenRow.FamilyID, &tokenRow.RevokedAt, &tokenRow.ReplacedByID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &RefreshToken{
		ID:        tokenRow.ID,
		UserID:    tokenRow.UserID,
		TokenHash: tokenRow.TokenHash,
		ExpiresAt: tokenRow.ExpiresAt.Unix(),
		IPAddress: tokenRow.IPAddress,
		FamilyID:  tokenRow.FamilyID,
		IsRevoked: tokenRow.RevokedAt.Valid,
	}, nil
}

func (r *RefreshTokenRepo) RevokeRefreshTokensForFamily(
	ctx context.Context,
	db db.DBTX,
	familyID string,
) error {
	var query = `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE family_id = $1
	`
	_, err := db.ExecContext(ctx, query, familyID)
	return err
}

func (r *RefreshTokenRepo) SetReplacedBy(
	ctx context.Context,
	db db.DBTX,
	oldTokenID string,
	newTokenID string,
) error {
	var query = `
		UPDATE refresh_tokens
		SET replaced_by_id = $1
		WHERE id = $2
	`
	_, err := db.ExecContext(ctx, query, newTokenID, oldTokenID)
	return err
}
