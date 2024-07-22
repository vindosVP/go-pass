// Package jwt works with jwt-tokens
package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/vindosVP/go-pass/internal/models"
)

// NewToken creates new JWT token for given user.
func NewToken(user *models.User, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies users token.
func VerifyToken(tokenString string, secret string) (string, int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", 0, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", 0, fmt.Errorf("invalid token: %w", err)
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", 0, fmt.Errorf("failed to extract email from token")
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		return "", 0, fmt.Errorf("failed to extract uid from token")
	}
	return email, int(uid), nil
}
