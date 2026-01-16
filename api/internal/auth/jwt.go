package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey            string
	tokenLifetimeMinutes int
}

type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, tokenLifetimeMinutes int) *jwtService {
	return &jwtService{
		secretKey:            secretKey,
		tokenLifetimeMinutes: tokenLifetimeMinutes,
	}
}

func (s *jwtService) generateJWT(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(s.tokenLifetimeMinutes) * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) parseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(s.secretKey), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
