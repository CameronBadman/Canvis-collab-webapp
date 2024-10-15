package database

import (
	"github.com/go-redis/redis/v8"
)

func InitRedis(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return client, nil
}
