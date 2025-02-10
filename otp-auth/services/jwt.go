package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT - JWT Token Generate Karega
func GenerateJWT(mobile string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	// Debugging: Print the SECRET_KEY
	fmt.Println("SECRET_KEY used for signing:", secretKey)

	if secretKey == "" {
		return "", errors.New("SECRET_KEY is not set")
	}

	claims := jwt.MapClaims{
		"mobile": mobile,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ✅ ParseJWT - JWT Token Verify Aur Decode Karega
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// Remove "Bearer " if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
