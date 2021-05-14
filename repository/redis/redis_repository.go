package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type CacheRepository interface {
	SetKey(key string, value []byte, ttl int)
	Get(key string) []byte
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository() *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_URL"),
		Password: viper.GetString("REDIS_PASSWORD"),
	})

	return &RedisRepository{client: client}
}

func (repository *RedisRepository) SetKey(key string, value interface{}, ttl int) {
	byteData, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
		return
	}

	duration, _ := time.ParseDuration(strconv.FormatInt(int64(ttl), 10))
	status := repository.client.Set(context.TODO(), key, string(byteData), duration)
	_, err = status.Result()
	if err != nil {
		fmt.Println(err)
	}
}

func (repository *RedisRepository) Get(key string) []byte {
	status := repository.client.Get(context.TODO(), key)
	stringResult, err := status.Result()
	if err != nil {
		fmt.Println(err)
	}

	return []byte(stringResult)
}
