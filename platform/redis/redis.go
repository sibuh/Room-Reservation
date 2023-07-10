package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(c *redis.Client) *redisAdapter {
	return &redisAdapter{
		client: c,
	}
}
func (r *redisAdapter) SetToRedis(ctx context.Context, key string, value interface{}) error {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, byteValue, 1*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil

}
func (r *redisAdapter) GetFromRedis(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {

		return "", err
	}
	return val, nil
}
