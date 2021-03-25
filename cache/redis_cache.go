package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type redisCache struct {
	host     string
	password string
}

func NewRedisCache() Cache {
	return &redisCache{
		host:     viper.GetString("REDIS_URL"),
		password: viper.GetString("REDIS_PASSWORD"),
	}
}

func (c *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: c.password,
	})
}

func (c *redisCache) getKeyPrefix() string {
	return "messaging-service"
}

func (c *redisCache) Get(key string) interface{} {
	key = c.getKeyPrefix() + key

	client := c.getClient()

	var ctx = context.Background()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	return val
}

func (c *redisCache) Set(key string, value interface{}) {
	panic("implement me")
}
