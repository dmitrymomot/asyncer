package asyncer

import (
	"fmt"

	"github.com/hibiken/asynq"
)

// NewClient creates a new asynq client from the given redis client instance.
func NewClient(redisConnStr string) (*asynq.Client, asynq.RedisConnOpt, error) {
	// Redis connect options for asynq client
	redisConnOpt, err := asynq.ParseRedisURI(redisConnStr)
	if err != nil {
		return nil, nil, fmt.Errorf("worket.NewClient: failed to parse redis connection string: %w", err)
	}

	// Init asynq client
	return asynq.NewClient(redisConnOpt), redisConnOpt, nil
}
