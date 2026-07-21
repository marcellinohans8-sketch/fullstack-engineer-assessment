package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	_ = godotenv.Load()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := client.Ping(Ctx).Result(); err != nil {
		RedisClient = nil
		fmt.Println("Redis unavailable, cache disabled:", err)
		return
	}

	RedisClient = client
	fmt.Println("Redis connected")
}
