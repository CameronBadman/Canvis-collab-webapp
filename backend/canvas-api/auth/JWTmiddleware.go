package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// TokenData represents the structure of a token stored in Redis
type TokenData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// JWTMiddleware is a middleware for validating JWT tokens using the provided Redis client
func JWTMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			// Log the full token for debugging (avoid in production)
			log.Printf("Authorization header: %s", authHeader)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(signingKey), nil
			})

			// Log token parsing results
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				log.Printf("Token validation error: %v", err)
				return
			}

			// Extract claims from the token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				log.Println("Invalid token claims")
				return
			}

			// Extract the user ID from the claims
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
			valid, err := validateTokenWithRedis(redisClient, userID, tokenString)
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
}

// validateTokenWithRedis validates the token by comparing it with the stored token in Redis
func validateTokenWithRedis(redisClient *redis.Client, userID, providedToken string) (bool, error) {
	// Construct the Redis key
	redisKey := "token:" + userID

	// Retrieve the stored token data from Redis
	tokenJson, err := redisClient.Get(context.Background(), redisKey).Result()
	if err == redis.Nil {
		return false, errors.New("token not found")
	} else if err != nil {
		return false, err
	}

	// Unmarshal the token data into TokenData struct
	var tokenData TokenData
	err = json.Unmarshal([]byte(tokenJson), &tokenData)
	if err != nil {
		return false, err
	}

	// Check if the token is expired
	if time.Now().Unix() > tokenData.ExpiresAt {
		return false, errors.New("token is expired")
	}

	// Check if the provided token matches the stored token
	if tokenData.AccessToken != providedToken {
		return false, errors.New("token mismatch")
	}

	return true, nil
}
