package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Username: "",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis ping error: %v", err)
	}
	return rdb
}
