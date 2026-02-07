package auth

import (
	"bs-books-api/internal/auth/refresh_token"
	"bs-books-api/internal/db"
	"bs-books-api/internal/users"
	"context"
)

type AuthService struct {
	db                  db.DBTX
	userService         *users.UserService
	jwtService          *JWTService
	refreshTokenService *refresh_token.RefreshTokenService
}

func NewAuthService(db db.DBTX, userService *users.UserService, jwtService *JWTService, refreshTokenService *refresh_token.RefreshTokenService) *AuthService {
	return &AuthService{
		db:                  db,
		userService:         userService,
		jwtService:          jwtService,
		refreshTokenService: refreshTokenService,
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

func (s *AuthService) Login(ctx context.Context, email, password string, ipAddress string) (string, *refresh_token.RefreshToken, error) {
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

	refreshToken, err := s.refreshTokenService.NewSession(ctx, user.ID, ipAddress)

	if err != nil {
		return "", nil, err
	}

	return token, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, oldRefreshToken string, ipAddress string) (string, *refresh_token.RefreshToken, error) {

	return "", nil, nil
}
