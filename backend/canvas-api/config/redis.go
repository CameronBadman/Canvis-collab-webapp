package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

// RedisCtx is a shared context for Redis operations
var RedisCtx = context.Background()

// InitRedis initializes a new Redis client with provided configuration
func InitRedis(redisHost, redisPort, redisPassword string) *redis.Client {
	// Construct Redis address
	redisAddr := redisHost + ":" + redisPort

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // Add if necessary
		DB:       0,             // Default DB
	})

	// Test Redis connection
	_, err := client.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
	}

	log.Printf("Connected to Redis at %s", redisAddr)
	return client
}
