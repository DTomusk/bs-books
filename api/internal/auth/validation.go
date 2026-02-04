package auth

import (
	"net/mail"
	"strings"
)

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrShortPassword
	}
	if len(password) > 64 {
		return ErrLongPassword
	}
	return nil
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

func validateUsername(username string) error {
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	if len(username) > 50 {
		return ErrUsernameTooLong
	}
	return nil
}
