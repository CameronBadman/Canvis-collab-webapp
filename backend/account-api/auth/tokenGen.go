package auth

import (
	"fmt"
	"time"

	"account-api/config"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT generates a JWT for the user after successful authentication
func GenerateJWT(userID string) (string, error) {
	// Define the claims for the JWT
	claims := jwt.MapClaims{
		"sub": userID,                                // Subject - the user ID
		"iat": time.Now().Unix(),                     // Issued at time
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Expiration time (24 hours)
	}

	// Use the JWT secret key for signing the token (HS256)
	secret := []byte(config.JWTSecretKey)

	// Create the JWT token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}
