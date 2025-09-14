package utils

import (
	"fmt"
	"os"
	"time"

	"auth-service/internal/model"

	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := JwtClaims{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "self-pickup",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}
	return token.SignedString([]byte(secret))
}

func GenerateAdminToken(admin *model.Admin) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := JwtClaims{
		Id:       admin.Id,
		Username: admin.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "self-pickup",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}
	return token.SignedString([]byte(secret))
}

func ValidateToken(signedToken string) (*JwtClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(signedToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
