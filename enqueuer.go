package asyncer

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
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

// NewEnqueuer creates a new email enqueuer.
// This function accepts EnqueuerOption to configure the enqueuer.
// Default values are used if no option is provided.
// Default values are:
//   - queue name: "default"
//   - task deadline: 1 minute
//   - max retry: 3
func NewEnqueuer(client *asynq.Client, opt ...EnqueuerOption) *Enqueuer {
	if client == nil {
		panic("client is nil")
	}

	e := &Enqueuer{
		client:       client,
		queueName:    "default",
		taskDeadline: time.Minute,
		maxRetry:     3,
	}

	for _, o := range opt {
		o(e)
	}

	return e
}

// EnqueueTask enqueues a task to the queue.
// This function returns an error if the task could not be enqueued.
// The task is enqueued with the following options:
//   - queue name: e.queueName
//   - task deadline: e.taskDeadline
//   - max retry: e.maxRetry
//   - unique: e.taskDeadline
func (e *Enqueuer) EnqueueTask(ctx context.Context, task *asynq.Task) error {
	if _, err := e.client.Enqueue(
		task,
		asynq.Queue(e.queueName),
		asynq.Deadline(time.Now().Add(e.taskDeadline)),
		asynq.MaxRetry(e.maxRetry),
		asynq.Unique(e.taskDeadline),
	); err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
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
