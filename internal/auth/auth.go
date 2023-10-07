package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthService interface {
	GenerateToken(userId int) (string, error)
	VerifyToken(token string) (*jwt.Token, error)
}

type authService struct{}

const (
	SECRET_KEY = "post12345"
)

func NewAuthService() AuthService {
	return &authService{}
}

func (a *authService) GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *authService) VerifyToken(token string) (*jwt.Token, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", fmt.Errorf("invalid method token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return tokenParsed, err
	}

	return tokenParsed, nil
}
