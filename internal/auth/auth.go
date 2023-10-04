package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

const (
	SECRET_KEY = "KEY-POST12345"
)

func (a *authService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenSigned, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenSigned, nil
}

func (a *authService) ValidateToken(token string) (*jwt.Token, error) {
	validToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return validToken, err
	}

	return validToken, nil
}
