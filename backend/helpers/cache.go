package helpers

import (
	"encoding/json"
	"time"

	"backend/config"
)

func GetCache(key string, dest interface{}) error {
	value, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), dest)
}

func SetCache(key string, value interface{}) error {
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
	keys, err := config.RedisClient.Keys(config.Ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return config.RedisClient.Del(config.Ctx, keys...).Err()
	}

	return nil
}