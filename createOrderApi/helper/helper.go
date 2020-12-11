package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client

func ConnectRd() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

	_, err = NewRedisHelper(redisClient)
	fmt.Println(err)
}

type RedisHelper interface {
	Set(key string, value interface{}, expireTime int64) error
	Get(key string) (string, error)
}

type redisHelper struct {
	clientSingleNode *redis.Client
}

//
func NewRedisHelper(client *redis.Client) (RedisHelper, error) {
	return &redisHelper{
		clientSingleNode: client,
	}, nil
}

func (s *redisHelper) Set(key string, value interface{}, expired int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = s.clientSingleNode.Set(key, data, time.Duration(expired)*time.Second).Result()
	return err
}

func (s *redisHelper) Get(key string) (string, error) {
	data, err := s.clientSingleNode.Get(key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
