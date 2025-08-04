package utils

import (
	"backend/dtos"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string) (*dtos.Claims, error) {
	claims := &dtos.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
