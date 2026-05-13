package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, role string, secret string, expiresIn time.Duration) (string, error) {
	claims := JWTCustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID.String(),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(expiresIn),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
