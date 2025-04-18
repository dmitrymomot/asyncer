package asyncer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type (
	// Enqueuer is a helper struct for enqueuing tasks.
	// You can encapsulate this struct in your own struct to add queue methods.
	// See pkg/worker/_example/enqueuer.go for an example.
	Enqueuer struct {
		client       *asynq.Client
		queueName    string
		taskDeadline time.Duration
		maxRetry     int
	}

	// EnqueuerOption is a function that configures an enqueuer.
	EnqueuerOption func(*Enqueuer)
)

// NewEnqueuerWithAsynqClient creates a new Enqueuer with the given Asynq client and options.
// It returns a pointer to the Enqueuer and an error if the Asynq client is nil.
// The Enqueuer is responsible for enqueueing tasks to the Asynq server.
// Default values are used if no option is provided.
// Default values are:
//   - queue name: "default"
//   - task deadline: 1 minute
//   - max retry: 3
func NewEnqueuerWithAsynqClient(client *asynq.Client, opt ...EnqueuerOption) (*Enqueuer, error) {
	if client == nil {
		return nil, ErrMissedAsynqClient
	}

	e := &Enqueuer{
		client:       client,
		queueName:    queueName,
		taskDeadline: time.Minute,
		maxRetry:     3,
	}

	for _, o := range opt {
		o(e)
	}

	return e, nil
}

// MustNewEnqueuerWithAsynqClient creates a new Enqueuer with the given Asynq client and options.
// It panics if an error occurs during the creation of the Enqueuer.
func MustNewEnqueuerWithAsynqClient(client *asynq.Client, opt ...EnqueuerOption) *Enqueuer {
	e, err := NewEnqueuerWithAsynqClient(client, opt...)
	if err != nil {
		panic(err)
	}

	return e
}

// NewEnqueuer creates a new Enqueuer with the given Redis connection string and options.
// Default values are used if no option is provided.
// It returns a pointer to the Enqueuer and an error if there was a problem creating the Enqueuer.
func NewEnqueuer(redisClient redis.UniversalClient, opt ...EnqueuerOption) (*Enqueuer, error) {
	client, err := NewClient(redisClient)
	if err != nil {
		return nil, errors.Join(ErrFailedToCreateEnqueuerWithClient, err)
	}

	return NewEnqueuerWithAsynqClient(client, opt...)
}

// MustNewEnqueuer creates a new Enqueuer with the given Redis connection string and options.
// It panics if an error occurs during the creation of the Enqueuer.
func MustNewEnqueuer(redisClient redis.UniversalClient, opt ...EnqueuerOption) *Enqueuer {
	e, err := NewEnqueuer(redisClient, opt...)
	if err != nil {
		panic(err)
	}

	return e
}

// EnqueueTask enqueues a task to be processed asynchronously.
// It takes a context and a task as parameters.
// The task is enqueued with the specified queue name, deadline, maximum retry count, and uniqueness constraint.
// Returns an error if the task fails to enqueue.
func (e *Enqueuer) EnqueueTask(ctx context.Context, taskName string, payload any, opts ...TaskOption) error {
	// Marshal payload to JSON bytes
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Join(ErrFailedToEnqueueTask, err)
	}

	// Set default options for enqueuing task.
	// These options can be overridden by the user provided options.
	defaultOptions := []asynq.Option{
		asynq.Queue(e.queueName),
		asynq.Deadline(time.Now().Add(e.taskDeadline)),
		asynq.MaxRetry(e.maxRetry),
		asynq.Unique(e.taskDeadline),
	}

	// Enqueue task
	if _, err := e.client.Enqueue(
		asynq.NewTask(taskName, jsonPayload),
		append(defaultOptions, opts...)...,
	); err != nil {
		return errors.Join(ErrFailedToEnqueueTask, err)
	}

	return nil
}

// Close closes the Enqueuer and releases any resources associated with it.
// It returns an error if there was a problem closing the Enqueuer.
func (e *Enqueuer) Close() error {
	if err := e.client.Close(); err != nil {
		return errors.Join(ErrFailedToCloseEnqueuer, err)
	}

	return nil
}

// WithQueueNameEnq configures the queue name for enqueuing.
// The queue name is the name of the queue where the task will be enqueued.
func WithQueueNameEnq(name string) EnqueuerOption {
	return func(e *Enqueuer) {
		e.queueName = name
	}
}

// WithTaskDeadline configures the task deadline.
// The task deadline is the time limit for the task to be processed.
func WithTaskDeadline(d time.Duration) EnqueuerOption {
	return func(e *Enqueuer) {
		e.taskDeadline = d
	}
}

// WithMaxRetry configures the max retry.
// The max retry is the number of times the task will be retried if it fails.
func WithMaxRetry(n int) EnqueuerOption {
	return func(e *Enqueuer) {
		e.maxRetry = n
	}
}
