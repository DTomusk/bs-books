package auth

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/users"
	"context"
)

type AuthService struct {
	db          db.DBTX
	userService *users.UserService
}

func NewAuthService(db db.DBTX, userService *users.UserService) *AuthService {
	return &AuthService{
		db:          db,
		userService: userService,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	existing_user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return err
	}

	if existing_user != nil {
		return ErrEmailAlreadyInUse
	}

	password_hash, err := hashPassword(password)

	if err != nil {
		return err
	}

	err = s.userService.CreateUser(email, password_hash, ctx)

	return err
}

func (s *AuthService) Login(ctx context.Context, email, password string) error {
	user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrInvalidCredentials
	}

	if err := comparePassword(user.PasswordHash, password); err != nil {
		return ErrInvalidCredentials
	}

	// TODO: generate and return JWT token

	return nil
}
