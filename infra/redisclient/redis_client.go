package redisclient

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

type RedisClient interface {
	SubscribeChannel(chatId string) *redis.PubSub
	PublishMessage(chatId string, msg string)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient() RedisClient {
	rdb := redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     viper.GetString("REDIS_URL"),
			Password: viper.GetString("REDIS_PASSWORD"),
		}),
	}
	pong, err := rdb.client.Ping(context.TODO()).Result()
	color.Green(pong, err)
	if err != nil {
		log.Fatal(err)
	}
	return &rdb
}

func (rc *redisClient) SubscribeChannel(chatId string) *redis.PubSub {
	return rc.client.Subscribe(context.TODO(), chatId)
}

func (rc *redisClient) PublishMessage(chatId string, msg string) {
	err := rc.client.Publish(context.TODO(), chatId, msg).Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}
