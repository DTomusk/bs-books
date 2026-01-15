package auth

import "fmt"

var (
	ErrShortPassword     = fmt.Errorf("password must be at least 8 characters long")
	ErrLongPassword      = fmt.Errorf("password must not exceed 64 characters")
	ErrInvalidEmail      = fmt.Errorf("invalid email address")
	ErrEmailAlreadyInUse = fmt.Errorf("email already in use")
)
