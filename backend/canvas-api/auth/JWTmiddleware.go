package auth

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"canvas-api/config"
)

// TokenData represents the structure of a token stored in Redis
type TokenData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// JWTMiddleware is a middleware for validating JWT tokens
func JWTMiddleware(next http.Handler) http.Handler {
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

		// Retrieve the JWT signing key from environment variables
		signingKey := os.Getenv("JWT_SECRET_KEY")
		if signingKey == "" {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("JWT_SECRET_KEY environment variable is not set")
			return
		}

		// Log the full token for debugging (beware of logging sensitive data in production)
		log.Printf("Authorization header: %s", authHeader)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(signingKey), nil
		})

		log.Printf("Authorization Token: %s", token.Raw)

		log.Printf("token valid: %t", !token.Valid)

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			log.Println("Invalid token claims")
			return
		}

		// Log claims for debugging
		log.Printf("Token Claims: %+v", claims)

		// Extract the user ID from the claims (assuming it's in the 'sub' field)
		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "User ID missing in token claims", http.StatusUnauthorized)
			log.Println("User ID missing in token claims")
			return
		}

		// Check for token expiration
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				log.Println("Token has expired")
				return
			}
		} else {
			http.Error(w, "Missing exp claim in token", http.StatusUnauthorized)
			log.Println("Missing exp claim in token")
			return
		}

		// Validate the token with Redis
		valid, err := validateTokenWithRedis(userID, tokenString)
		if err != nil || !valid {
			http.Error(w, "Token validation failed", http.StatusUnauthorized)
			log.Printf("Token validation failed for user %s: %v", userID, err)
			return
		}

		// Log successful validation
		log.Printf("Token validated successfully for user %s", userID)

		// Pass the request to the next handler

		ctx := SetUserIDInContext(r.Context(), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateTokenWithRedis validates the token by comparing it with the stored token in Redis
func validateTokenWithRedis(userID, providedToken string) (bool, error) {
	// Construct the Redis key
	redisKey := "token:" + userID

	// Retrieve the stored token data from Redis
	tokenJson, err := config.RedisClient.Get(config.RedisCtx, redisKey).Result()
	if err != nil {
		//log.Printf("Failed to retrieve token from Redis for user %s: %v", userID, err)
		return false, errors.New("token not found or error retrieving token")
	}

	// Log the Redis data (token)
	//log.Printf("Retrieved token from Redis for user %s: %s", userID, tokenJson)

	// Unmarshal the token data into TokenData struct
	var tokenData TokenData
	err = json.Unmarshal([]byte(tokenJson), &tokenData)
	if err != nil {
		//log.Printf("Failed to unmarshal token data for user %s: %v", userID, err)
		return false, err
	}

	// Log the token data
	//log.Printf("Redis token data for user %s: %+v", userID, tokenData)

	// Check if the token is expired
	if time.Now().Unix() > tokenData.ExpiresAt {
		//log.Printf("Token for user %s is expired", userID)
		return false, errors.New("token is expired")
	}

	// Log token expiration check
	//log.Printf("Token expiration check passed for user %s", userID)

	// Check if the provided token matches the stored token
	if tokenData.AccessToken != providedToken {
		//log.Printf("Provided token for user %s does not match stored token", userID)
		//log.Printf("Provided token: %s", providedToken)       // Log provided token for comparison
		//log.Printf("Stored token: %s", tokenData.AccessToken) // Log stored token for comparison
		return false, errors.New("token mismatch")
	}

	// Token is valid
	//log.Printf("Token validation succeeded for user %s", userID)
	return true, nil
}
