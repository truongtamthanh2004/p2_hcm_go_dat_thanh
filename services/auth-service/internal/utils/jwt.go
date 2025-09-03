package utils

import (
	"auth-service/internal/model"
	"errors"
	"log"
	"os"
	"time"

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

func GenerateAccessToken(user *model.AuthUser) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID:     user.UserID,
		Email:      user.Email,
		Role:       user.Role,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

func GenerateRefreshToken(user *model.AuthUser) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID:     user.UserID,
		Email:      user.Email,
		Role:       user.Role,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
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
