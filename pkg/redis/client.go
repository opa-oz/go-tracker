package redis

import (
	"tracker/pkg/config"

	"github.com/redis/go-redis/v9"
)

func GetClient(cfg *config.Environment) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	return rdb
}
