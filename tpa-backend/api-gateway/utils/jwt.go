package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	AccessTokenExpiry  = time.Minute * 15  // Access token 15 minutes
	RefreshTokenExpiry = time.Hour * 24 * 7 // Refresh token 7 days
)

func CreateJWT(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(), // token expires in 24h
	})

	jwtKey := GetDotEnv("JWT_KEY")

	return token.SignedString([]byte(jwtKey))
}

func ParseJWT(tokenString string) (string, error) {
	jwtKey := GetDotEnv("JWT_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", fmt.Errorf("user_id not found or not a string")
	}

	return "", fmt.Errorf("invalid token claims")
}