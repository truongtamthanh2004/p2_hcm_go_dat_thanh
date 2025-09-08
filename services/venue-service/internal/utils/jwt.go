package utils

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`
	jwt.RegisteredClaims
}

var (
	jwtSecretKey string
	jwtIssuer    string
)

func InitJWT() {
	jwtSecretKey = os.Getenv("SECRET_KEY")
	jwtIssuer = os.Getenv("JWT_ISSUER")
	if jwtSecretKey == "" || jwtIssuer == "" {
		log.Fatal("Missing SECRET_KEY or JWT_ISSUER in environment")
	}
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error.unexpected_signing_method")
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("error.invalid_token")
	}
	return claims, nil
}
