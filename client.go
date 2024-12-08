package asyncer

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// NewClient creates a new instance of the asynq client using the provided Redis connection string.
// It returns the created client, the Redis connection options, and any error encountered during the process.
func NewClient(redisClient redis.UniversalClient) (*asynq.Client, error) {
	// Init asynq client
	return asynq.NewClientFromRedisClient(redisClient), nil
}
