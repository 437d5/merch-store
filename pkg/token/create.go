package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTExpAt = 1
)

type TokenClaims struct {
	Id       int
	Username string
	jwt.RegisteredClaims
}

func CreateToken(id int, username, secret string, expAt int) (string, error) {
	now := time.Now()
	exp := now.Add(time.Duration(expAt) * time.Hour)

	claims := TokenClaims{
		Id:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed token signing: %s", err)
	}

	return tokenString, nil
}
