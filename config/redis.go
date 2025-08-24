package config

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var redisInstance *redis.Client

func GetRedis() *redis.Client {
	if redisInstance == nil {
		redisInstance = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_URL"),
		})
	}
	return redisInstance
}

func GetPubSubClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
}
