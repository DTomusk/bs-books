package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtService struct {
	secretKey            string
	tokenLifetimeMinutes int
}

func NewJWTService(secretKey string, tokenLifetimeMinutes int) *jwtService {
	return &jwtService{
		secretKey:            secretKey,
		tokenLifetimeMinutes: tokenLifetimeMinutes,
	}
}

func (s *jwtService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": jwt.TimeFunc().Add(time.Duration(s.tokenLifetimeMinutes) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}
