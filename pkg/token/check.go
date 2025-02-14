package token

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid JWT")
)

func ValidateToken(token, secret string, id int) (bool, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return false, fmt.Errorf("invalid token: %s", err)
	}

	if !parsedToken.Valid {
		return false, ErrInvalidToken
	}

	return true, nil
}