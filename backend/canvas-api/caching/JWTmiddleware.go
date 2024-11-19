package caching

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
)

// JWTMiddleware is a middleware for validating JWT tokens
func JWTMiddleware(next http.Handler, requiredUserID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			log.Println("Authorization header is missing")
			return
		}

		// Parse the token from the "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token format is invalid", http.StatusUnauthorized)
			log.Println("Bearer token format is invalid")
			return
		}

		// Retrieve the stored token from Redis
		storedTokenData, err := GetToken(requiredUserID)
		if err != nil {
			http.Error(w, "Failed to retrieve token from Redis", http.StatusInternalServerError)
			log.Printf("Error retrieving token for user %s from Redis: %v", requiredUserID, err)
			return
		}

		// Check if the token data was found and if the token matches
		if storedTokenData == nil || storedTokenData.AccessToken != tokenString {
			http.Error(w, "Token mismatch or expired", http.StatusUnauthorized)
			log.Printf("Token mismatch for user %s", requiredUserID)
			return
		}

		// Validate the JWT token
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Here you would validate the token's signature with your secret or public key
			return nil, nil
		})
		if err != nil || !parsedToken.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			log.Printf("Invalid or expired token: %v", err)
			return
		}

		// Extract claims from the token
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			log.Println("Invalid token claims")
			return
		}

		// Extract userID from the token claims (assume it's in the 'sub' field)
		userID, ok := claims["sub"].(string)
		if !ok || userID != requiredUserID {
			http.Error(w, "User ID mismatch", http.StatusUnauthorized)
			log.Printf("User ID mismatch: expected %s but got %s", requiredUserID, userID)
			return
		}

		// Log successful validation
		log.Printf("User ID %s matches token, proceeding to next handler", userID)

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
