package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Redis struct {
	client *redis.Client
}

func NewClient() *Redis {
	return &Redis{redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})}
}

func (redis *Redis) Set(key string, val interface{}) error {
	err := redis.client.Set(ctx, key, val, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (redis *Redis) Get(key string) interface{} {
	val, err := redis.client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(key, val)
	return val
}

func (redis *Redis) Del(key string) error {
	err := redis.client.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (redis *Redis) AddToSet(key string, val interface{}) error {
	err := redis.client.SAdd(ctx, key, val).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (redis *Redis) GetRandomFromSet(key string) interface{} {
	val, err := redis.client.SRandMember(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(key, val)
	return val
}

func (redis *Redis) DeleteAll(key string) error {
	err := redis.client.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
