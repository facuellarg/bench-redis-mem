package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

func (r *Redis) Set(key string, value string) error {
	return r.client.Set(r.client.Context(), key, value, 0).Err()
}

func (r *Redis) Clean() error {
	return r.client.FlushAll(r.client.Context()).Err()
}
