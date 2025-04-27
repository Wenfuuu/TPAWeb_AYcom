package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenExpiry  = time.Minute * 15  // Access token  15 minutes
	RefreshTokenExpiry = time.Hour * 24 * 7 // Refresh token 7 days
)

func CreateJWT(userID string, userRole string) (string, string, error) {
	jwtKey := GetDotEnv("JWT_KEY")
	refreshKey := GetDotEnv("JWT_REFRESH_KEY")

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"user_role": userRole,
		"exp":     time.Now().Add(AccessTokenExpiry).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(RefreshTokenExpiry).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func parseToken(tokenString string, key string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, okID := claims["user_id"].(string)
		userRole, okRole := claims["user_role"].(string)

		if !okID {
			return "", "", fmt.Errorf("user_id not found or not a string")
		}
		if !okRole {
			return "", "", fmt.Errorf("user_role not found or not a string")
		}

		return userID, userRole, nil
	}

	return "", "", fmt.Errorf("invalid token claims")
}

func ParseAccessToken(tokenString string) (string, string, error) {
	jwtKey := GetDotEnv("JWT_KEY")

	return parseToken(tokenString, jwtKey)
}

func ParseRefreshToken(tokenString string) (string, string, error) {
	refreshKey := GetDotEnv("JWT_REFRESH_KEY")

	return parseToken(tokenString, refreshKey)
}