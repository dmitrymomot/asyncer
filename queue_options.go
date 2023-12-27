package asyncer

import (
	"time"

	"github.com/hibiken/asynq"
)

// WithQueues sets the queues.
// It panics if the sum of concurrency is not equal to the concurrency
// set in the config.
// If you want to increase the concurrency of a queue, you can use
// asyncer.WithQueueConcurrency before this option.
func WithQueues(queues map[string]int) QueueServerOption {
	return func(cnf *asynq.Config) {
		// check if sum of concurrency is equal to the concurrency
		// set in the config.
		var sum int
		for _, v := range queues {
			sum += v
		}
		if sum != cnf.Concurrency {
			panic("sum of concurrency is not equal to the concurrency set in the config")
		}
		cnf.Queues = queues
	}
}

// WithQueueName sets the queue name.
func WithQueueName(name string) QueueServerOption {
	return func(cnf *asynq.Config) {
		cnf.Queues = map[string]int{
			name: cnf.Concurrency,
		}
	}
}

// WithQueueConcurrency sets the queue concurrency.
func WithQueueConcurrency(concurrency int) QueueServerOption {
	return func(cnf *asynq.Config) {
		cnf.Concurrency = concurrency
		for k := range cnf.Queues {
			cnf.Queues[k] = concurrency
		}
	}
}

// WithQueueShutdownTimeout sets the queue shutdown timeout.
func WithQueueShutdownTimeout(timeout time.Duration) QueueServerOption {
	return func(cnf *asynq.Config) {
		cnf.ShutdownTimeout = timeout
	}
}

// WithQueueLogLevel sets the queue log level.
func WithQueueLogLevel(level string) QueueServerOption {
	return func(cnf *asynq.Config) {
		cnf.LogLevel = getAsynqLogLevel(level)
	}
}

// WithQueueLogger sets the queue logger.
func WithQueueLogger(logger asynq.Logger) QueueServerOption {
	return func(cnf *asynq.Config) {
		if logger != nil {
			cnf.Logger = logger
		}
	}
}
