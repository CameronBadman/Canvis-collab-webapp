package caching

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"account-api/config"
)

// TokenData represents the data to store in Redis for an access token.
type TokenData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// StoreToken stores the access token and its expiration time in Redis.
func StoreToken(userID, accessToken string, expiresIn int64) error {
	// Create token data with expiration time
	tokenData := TokenData{
		UserID:      userID,
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Unix() + expiresIn,
	}

	// Marshal the token data to JSON
	tokenJson, err := json.Marshal(tokenData)
	if err != nil {
		log.Printf("Failed to marshal token data: %v", err)
		return err
	}

	// Store in Redis with a key of token:user_id
	redisKey := "token:" + userID
	err = config.RedisClient.Set(config.RedisCtx, redisKey, tokenJson, time.Duration(expiresIn)*time.Second).Err()
	if err != nil {
		log.Printf("Failed to store token in Redis: %v", err)
		return err
	}

	log.Printf("Stored token for user %s in Redis with key %s", userID, redisKey)
	return nil
}

// CheckToken validates if the provided access token matches the stored token in Redis for the given user.
func CheckToken(providedToken, userID string) (bool, error) {
	// Construct the Redis key
	redisKey := "token:" + userID

	// Retrieve the stored token data from Redis
	tokenJson, err := config.RedisClient.Get(config.RedisCtx, redisKey).Result()
	if err != nil {
		log.Printf("Failed to retrieve token from Redis: %v", err)
		return false, errors.New("token not found or error retrieving token")
	}

	// Unmarshal the token JSON data
	var tokenData TokenData
	err = json.Unmarshal([]byte(tokenJson), &tokenData)
	if err != nil {
		log.Printf("Failed to unmarshal token data: %v", err)
		return false, err
	}

	// Check if the token has expired
	if time.Now().Unix() > tokenData.ExpiresAt {
		log.Printf("Token for user %s is expired", userID)
		return false, errors.New("token is expired")
	}

	// Check if the provided token matches the stored token
	if tokenData.AccessToken != providedToken {
		log.Printf("Provided token for user %s does not match", userID)
		return false, errors.New("provided token does not match stored token")
	}

	// Token is valid
	log.Printf("Token for user %s is valid", userID)
	return true, nil
}

// RetrieveJWTFromCache retrieves the JWT from Redis based on the userID
func RetrieveJWTFromCache(userID string) (string, error) {
	// Construct the Redis key and retrieve token data
	redisKey := "token:" + userID
	tokenJson, err := config.RedisClient.Get(config.RedisCtx, redisKey).Result()
	if err != nil {
		log.Printf("Failed to retrieve token from Redis for user %s: %v", userID, err)
		return "", errors.New("token not found or error retrieving token")
	}

	// Unmarshal token data from JSON
	var tokenData TokenData
	err = json.Unmarshal([]byte(tokenJson), &tokenData)
	if err != nil {
		log.Printf("Failed to unmarshal token data: %v", err)
		return "", err
	}

	// Check if the token has expired
	if time.Now().Unix() > tokenData.ExpiresAt {
		log.Printf("Token for user %s is expired", userID)
		return "", errors.New("token is expired")
	}

	return tokenData.AccessToken, nil
}
