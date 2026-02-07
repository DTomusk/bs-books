package auth

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/users"
	"context"
)

type AuthService struct {
	db                     db.DBTX
	repo                   *AuthRepo
	userService            *users.UserService
	jwtService             *JWTService
	refreshTokenExpiryDays int
}

func NewAuthService(db db.DBTX, repo *AuthRepo, userService *users.UserService, jwtService *JWTService, refreshTokenExpiryDays int) *AuthService {
	return &AuthService{
		db:                     db,
		repo:                   repo,
		userService:            userService,
		jwtService:             jwtService,
		refreshTokenExpiryDays: refreshTokenExpiryDays,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	if err := validateUsername(username); err != nil {
		return err
	}

	existing_user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return err
	}

	if existing_user != nil {
		return ErrEmailAlreadyInUse
	}

	existing_user, err = s.userService.GetUserByUsername(username, ctx)
	if err != nil {
		return err
	}

	if existing_user != nil {
		return ErrUsernameAlreadyInUse
	}

	password_hash, err := hashPassword(password)

	if err != nil {
		return err
	}

	err = s.userService.CreateUser(username, email, password_hash, ctx)

	return err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *RefreshToken, error) {
	user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := comparePassword(user.PasswordHash, password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	// TODO: revoke existing refresh tokens for user
	err = s.repo.RevokeRefreshTokensForUser(user.ID)

	if err != nil {
		return "", nil, err
	}

	refreshToken, err := NewRefreshToken(user.ID, s.refreshTokenExpiryDays)

	if err != nil {
		return "", nil, err
	}

	err = s.repo.SaveRefreshToken(refreshToken)

	if err != nil {
		return "", nil, err
	}

	return token, refreshToken, nil
}
