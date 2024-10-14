package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// AuthMiddleware checks for authorization using a token stored in Redis.
func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := tokenParts[1]
		firebaseUID := c.Request.Header.Get("X-Firebase-UID")

		// Check if the token exists in Redis for the given Firebase UID
		storedToken, err := redisClient.Get(firebaseUID).Result()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error verifying token"})
			return
		}

		if token != storedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// ForwardHeaders forwards specific headers from the incoming request to the response.
func ForwardHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Forward the Authorization header
		if auth := c.Request.Header.Get("Authorization"); auth != "" {
			c.Writer.Header().Set("Authorization", auth)
		}

		// Forward the X-Firebase-UID header
		if firebaseUID := c.Request.Header.Get("X-Firebase-UID"); firebaseUID != "" {
			c.Writer.Header().Set("X-Firebase-UID", firebaseUID)
		}

		c.Next()
	}
}
