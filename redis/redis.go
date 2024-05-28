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
		Password: "", // no password set
		DB:       0,  // use default DB
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
