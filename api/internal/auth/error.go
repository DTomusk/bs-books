package auth

import "fmt"

var (
	ErrShortPassword        = fmt.Errorf("password must be at least 8 characters long")
	ErrLongPassword         = fmt.Errorf("password must not exceed 64 characters")
	ErrInvalidEmail         = fmt.Errorf("invalid email address")
	ErrEmailAlreadyInUse    = fmt.Errorf("email already in use")
	ErrInvalidCredentials   = fmt.Errorf("invalid email or password")
	ErrInvalidToken         = fmt.Errorf("invalid token")
	ErrMissingAuthHeader    = fmt.Errorf("missing authorization header")
	ErrInvalidAuthHeader    = fmt.Errorf("invalid authorization header format")
	ErrUsernameTooShort     = fmt.Errorf("username must be at least 3 characters long")
	ErrUsernameTooLong      = fmt.Errorf("username must not exceed 50 characters")
	ErrInvalidUsername      = fmt.Errorf("username can only contain letters, numbers, underscores, or hyphens")
	ErrUsernameAlreadyInUse = fmt.Errorf("username already in use")
)
