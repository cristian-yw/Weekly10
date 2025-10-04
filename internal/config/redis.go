package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379" // default
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Redis ping error: %v", err)
	}

	return rdb
}
