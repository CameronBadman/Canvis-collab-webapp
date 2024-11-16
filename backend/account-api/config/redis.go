package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password for default setup
		DB:       0,                // Use default DB
	})

	// Test Redis connection
	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}
