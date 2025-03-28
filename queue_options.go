package asyncer

import (
	"time"

	"github.com/hibiken/asynq"
)

// WithQueues sets the queues with their priorities.
// The map key is the queue name and the value is the priority.
// Higher priority values give the queue higher processing preference.
func WithQueues(queues map[string]int) QueueServerOption {
	return func(cnf *asynq.Config) {
		if len(queues) > 0 {
			cnf.Queues = queues
		}
	}
}

// WithQueue sets the queue name.
func WithQueue(name string, priority int) QueueServerOption {
	return func(cnf *asynq.Config) {
		if priority < 1 {
			priority = 1
		}
		cnf.Queues = map[string]int{
			name: priority,
		}
	}
}

// WithQueueConcurrency sets the queue concurrency.
func WithQueueConcurrency(concurrency int) QueueServerOption {
	return func(cnf *asynq.Config) {
		if concurrency < 1 {
			concurrency = 1
		}
		cnf.Concurrency = concurrency
	}
}

// WithQueueShutdownTimeout sets the queue shutdown timeout.
func WithQueueShutdownTimeout(timeout time.Duration) QueueServerOption {
	return func(cnf *asynq.Config) {
		if timeout < 0 {
			timeout = 0
		}
		cnf.ShutdownTimeout = timeout
	}
}

// WithQueueLogLevel sets the queue log level.
func WithQueueLogLevel(level string) QueueServerOption {
	return func(cnf *asynq.Config) {
		cnf.LogLevel = castToAsynqLogLevel(level)
	}
}

// WithQueueStrictPriority sets the queue strict priority.
func WithQueueStrictPriority(strict bool) QueueServerOption {
	return func(cnf *asynq.Config) {
		cnf.StrictPriority = strict
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

// WithQueueErrorHandler sets the queue error handler.
func WithQueueErrorHandler(handler asynq.ErrorHandler) QueueServerOption {
	return func(cnf *asynq.Config) {
		if handler != nil {
			cnf.ErrorHandler = handler
		}
	}
}
