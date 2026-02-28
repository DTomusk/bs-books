package users

import "github.com/google/uuid"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Username     string
}

func NewUser(email, passwordHash, username string) *User {
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	}
}
