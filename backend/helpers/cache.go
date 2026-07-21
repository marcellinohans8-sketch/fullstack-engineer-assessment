package helpers

import (
	"encoding/json"
	"time"

	"backend/config"

	"github.com/redis/go-redis/v9"
)

func GetCache(key string, dest interface{}) error {
	if config.RedisClient == nil {
		return redis.Nil
	}

	value, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), dest)
}

func SetCache(key string, value interface{}) error {
	if config.RedisClient == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(
		config.Ctx,
		key,
		data,
		60*time.Second,
	).Err()
}

func DeleteCache(pattern string) error {
	if config.RedisClient == nil {
		return nil
	}

	var cursor uint64

	for {
		keys, nextCursor, err := config.RedisClient.Scan(config.Ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := config.RedisClient.Del(config.Ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	return nil
}
