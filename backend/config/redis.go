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
	Ctx = context.Background()
)

func ConnectRedis() {
	_ = godotenv.Load()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Redis Connected")
}