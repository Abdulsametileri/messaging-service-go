package redisclient

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisClient interface{}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient() RedisClient {
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     viper.GetString("REDIS_URL"),
			Password: viper.GetString("REDIS_PASSWORD"),
		}),
	}
}
