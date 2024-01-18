package asyncer

import (
	"errors"

	"github.com/hibiken/asynq"
)

// NewClient creates a new instance of the asynq client using the provided Redis connection string.
// It returns the created client, the Redis connection options, and any error encountered during the process.
func NewClient(redisConnStr string) (*asynq.Client, asynq.RedisConnOpt, error) {
	// Redis connect options for asynq client
	redisConnOpt, err := asynq.ParseRedisURI(redisConnStr)
	if err != nil {
		return nil, nil, errors.Join(ErrFailedToParseRedisURI, err)
	}

	// Init asynq client
	return asynq.NewClient(redisConnOpt), redisConnOpt, nil
}
