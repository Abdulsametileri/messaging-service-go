package redisservice

import (
	"github.com/Abdulsametileri/messaging-service/infra/redisclient"
	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	SubscribeChannel(chatId string) *redis.PubSub
	PublishMessage(chatId string, msg string)
}

type redisService struct {
	client redisclient.RedisClient
}

func NewRedisService(client redisclient.RedisClient) RedisService {
	return &redisService{client: client}
}

func (r *redisService) SubscribeChannel(chatId string) *redis.PubSub {
	return r.client.SubscribeChannel(chatId)
}

func (r *redisService) PublishMessage(chatId string, msg string) {
	r.client.PublishMessage(chatId, msg)
}
