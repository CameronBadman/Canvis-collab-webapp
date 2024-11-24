package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
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

// JWTMiddleware validates JWT tokens using Redis
func JWTMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse Authorization header
			authHeader := r.Header.Get("Authorization")
			log.Printf("Received Authorization header: %s", authHeader)

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
				log.Println("Authorization header is missing or invalid")
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			log.Printf("Extracted token: %s", tokenString)

			// Retrieve JWT signing key
			signingKey := os.Getenv("JWT_SECRET_KEY")
			if signingKey == "" {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Println("Missing JWT_SECRET_KEY environment variable")
				return
			}
			log.Println("Successfully retrieved JWT signing key from environment variables")

			// Parse and validate JWT
			claims, err := parseAndValidateJWT(tokenString, signingKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				log.Printf("JWT validation failed: %v", err)
				return
			}

			// Extract user ID from claims
			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				//http.Error(w, "Invalid token: missing user ID", http.StatusUnauthorized)
				log.Println("User ID missing in token claims")
				return
			}
			log.Printf("Extracted user ID: %s", userID)

			// Validate token against Redis
			valid, err := validateTokenWithRedis(redisClient, userID, tokenString)
			if err != nil || !valid {
				http.Error(w, "Token validation failed", http.StatusUnauthorized)
				log.Printf("Redis token validation failed for user %s: %v", userID, err)
				return
			}
			log.Printf("Token validated successfully for user %s", userID)

			// Set user ID in context and proceed
			ctx := SetUserIDInContext(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// parseAndValidateJWT parses the JWT and validates its signature and expiration
func parseAndValidateJWT(tokenString, signingKey string) (jwt.MapClaims, error) {
	log.Printf("Parsing JWT token: %s", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Method)
			return nil, errors.New("unexpected signing method")
		}
		log.Println("Successfully verified signing method: HMAC")
		return []byte(signingKey), nil
	})
	if err != nil || !token.Valid {
		log.Printf("JWT token is invalid: %v", err)
		//return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid token claims")
		return nil, errors.New("invalid token claims")
	}

	// Validate expiration
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		log.Println("Token has expired or invalid exp claim")
		return nil, errors.New("token has expired")
	}

	log.Printf("Token is valid with exp: %v", exp)
	return claims, nil
}

// validateTokenWithRedis validates the token against the stored token in Redis
func validateTokenWithRedis(redisClient *redis.Client, userID, providedToken string) (bool, error) {
	redisKey := "token:" + userID // Use a configurable prefix if needed

	// Fetch token from Redis
	log.Printf("Fetching token from Redis for user: %s, key: %s", userID, redisKey)
	tokenJson, err := redisClient.Get(context.Background(), redisKey).Result()
	if err == redis.Nil {
		log.Printf("Token not found in Redis for user %s", userID)
		return false, errors.New("token not found in Redis")
	} else if err != nil {
		log.Printf("Redis error while fetching token for user %s: %v", userID, err)
		return false, err
	}

	// Unmarshal token data
	var tokenData TokenData
	if err := json.Unmarshal([]byte(tokenJson), &tokenData); err != nil {
		log.Printf("Failed to unmarshal token data for user %s: %v", userID, err)
		return false, errors.New("failed to parse token data")
	}

	// Check token expiration and match
	if time.Now().Unix() > tokenData.ExpiresAt {
		log.Printf("Token has expired for user %s", userID)
		return false, errors.New("token has expired")
	}
	if tokenData.AccessToken != providedToken {
		log.Printf("Token mismatch for user %s: expected %s, got %s", userID, tokenData.AccessToken, providedToken)
		return false, errors.New("token mismatch")
	}

	log.Printf("Token validated successfully for user %s", userID)
	return true, nil
}
