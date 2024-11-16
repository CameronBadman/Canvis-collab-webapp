package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

// ExtractSubFromIDToken extracts the "sub" claim from a Cognito ID token
func ExtractSubFromIDToken(idToken string) (string, error) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse ID token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return "", fmt.Errorf("sub not found in ID token")
	}

	return claims["sub"].(string), nil
}
