package auth

import "golang.org/x/crypto/bcrypt"

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrShortPassword
	}
	if len(password) > 64 {
		return ErrLongPassword
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
