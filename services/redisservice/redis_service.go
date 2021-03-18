package redisservice

import "github.com/Abdulsametileri/messaging-service/infra/redisclient"

type RedisService interface {
}

type redisService struct {
	client redisclient.RedisClient
}

func NewRedisService(client redisclient.RedisClient) RedisService {
	return &redisService{client: client}
}
