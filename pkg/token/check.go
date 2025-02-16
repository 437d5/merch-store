package token

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid JWT")
)

func ValidateToken(token, secret string) (int, bool, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, false, fmt.Errorf("invalid token: %s", err)
	}

	if !parsedToken.Valid {
		return 0, false, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		return 0, false, ErrInvalidToken
	}

	return claims.Id, true, nil
}
