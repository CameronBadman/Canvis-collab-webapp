package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	// Load Redis configuration from environment variables
	redisAddr := os.Getenv("AUTH_REDIS_HOST") + ":" + os.Getenv("AUTH_REDIS_PASSWORD")
	redisPassword := os.Getenv("AUTH_REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis server address
		Password: redisPassword, // Redis password
		DB:       0,             // Use default DB
	})

	// Test Redis connection
	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}
