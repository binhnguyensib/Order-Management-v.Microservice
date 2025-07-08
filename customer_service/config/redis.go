package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	rdb    *redis.Client
	Logger = logrus.New()
)

func InitRedis() {
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		Logger.Errorf("Unable to connect to Redis: %v", err)
	}
	Logger.Infof("Connected to Redis")
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	cmd := rdb.Set(ctx, key, value, ttl)
	return cmd.Err()
}

func Get(ctx context.Context, key string) (string, error) {
	cmd := rdb.Get(ctx, key)
	return cmd.Val(), cmd.Err()
}

func Del(ctx context.Context, key string) error {
	cmd := rdb.Del(ctx, key)
	return cmd.Err()
}
