package caching

import (
	"canvas-api/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type TokenData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// GetToken retrieves the stored token and expiration time for a given user from Redis
func GetToken(userID string) (*TokenData, error) {
	ctx := context.Background()
	key := fmt.Sprintf("token:%s", userID)

	// Fetch token data from Redis
	tokenJson, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Token not found in Redis
			return nil, nil
		}
		// Redis error
		log.Printf("Error retrieving token for user %s from Redis: %v", userID, err)
		return nil, err
	}

	// Deserialize the token data from JSON
	var tokenData TokenData
	if err := json.Unmarshal([]byte(tokenJson), &tokenData); err != nil {
		log.Printf("Failed to unmarshal token data for user %s: %v", userID, err)
		return nil, err
	}

	// Check if the token has expired
	if time.Now().Unix() > tokenData.ExpiresAt {
		// Token has expired
		log.Printf("Token for user %s has expired", userID)
		return nil, nil
	}

	// Return the token data
	return &tokenData, nil
}
