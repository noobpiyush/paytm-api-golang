package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type customClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateToken(username string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	if len(secretKey) == 0 {
		return "", errors.New("JWT_SECRET env var not found")
	}

	claims := customClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 36)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", fmt.Errorf("failedd to sign token: %w", err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	if len(secretKey) == 0 {
		return "", errors.New("JWT_SECRET env variable not found")
	}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}

	return claims.Username, nil
}

func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}
