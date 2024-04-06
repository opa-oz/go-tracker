package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetRedis(c *gin.Context) (*redis.Client, error) {
	r := c.Value("Redis")

	if r == nil {
		err := fmt.Errorf("could not retrieve Redis")
		return nil, err
	}

	rdb, ok := r.(*redis.Client)
	if !ok {
		err := fmt.Errorf("variable Redis has wrong type")
		return nil, err
	}

	return rdb, nil
}
