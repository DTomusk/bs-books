package auth

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/users"
	"context"
	"net/mail"
	"strings"
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

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return ErrInvalidEmail
	}

	domain := parts[1]

	if !strings.Contains(domain, ".") {
		return ErrInvalidEmail
	}

	return nil
}
