package auth

type AuthRepo struct{}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{}
}

func (r *AuthRepo) RevokeRefreshTokensForUser(userID string) error {
	return nil
}

func (r *AuthRepo) SaveRefreshToken(token *RefreshToken) error {
	return nil
}
