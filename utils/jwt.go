// utils/jwt.go
package utils

import (
	"fmt"
	"go-clickhouse-example/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key for signing the JWT token (in a real app, use environment variables or a secure vault)
var secretKey = []byte("test_secret_key")

// GenerateJWT generates a JWT token for the authenticated user
func GenerateJWT(user *models.UserResponse) (string, error) {
	// Define the token claims with a Unix timestamp for expiration
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours (Unix timestamp)
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %v", err)
	}

	return tokenString, nil
}

// ParseJWT parses and validates the JWT token
func ParseJWT(tokenString string) (*models.User, error) {
	// Print the raw token to inspect its parts
	fmt.Println("Token:", tokenString)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse the token: %v", err)
	}

	// Extract claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Print the claims to check their structure
		fmt.Println("Claims:", claims)

		// Decode the expiration time (Unix timestamp)
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid 'exp' claim in token")
		}
		expirationTime := time.Unix(int64(exp), 0)

		// Check if the token has expired
		if time.Now().After(expirationTime) {
			return nil, fmt.Errorf("token is expired")
		}

		// Return the user information from the claims
		user := &models.User{
			ID:   uint64(claims["user_id"].(float64)),
			Role: claims["role"].(string),
		}
		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}
