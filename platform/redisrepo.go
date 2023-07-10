package platform

import "context"

type RedisInterface interface {
	SetToRedis(ctx context.Context, key string, value interface{}) error
	GetFromRedis(ctx context.Context, key string) (string, error)
}
