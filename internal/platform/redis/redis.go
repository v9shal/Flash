package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, redisURL string) (*redis.Client, error) {
	url := redisURL
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	err = client.Ping(ctx).Err()
	if err != nil {

		return nil, fmt.Errorf("error while pinging the redis %w", err)
	}
	return client, nil
}
