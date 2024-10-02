package asyncer

import (
	"context"

	"github.com/hibiken/asynq"
)

type (

	// TaskScheduler is an interface for task schedulers.
	// It is used to register task schedulers in the queue server.
	TaskScheduler interface {
		// TaskName returns the name of the task. It is used to register the task handler.
		TaskName() string
		// Schedule returns the cron spec for the task.
		// For more information about cron spec, see https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format.
		Schedule() string
		// Options returns the options for the task scheduler.
		Options() []TaskOption
	}

	// scheduledHandlerFunc is a function that handles a scheduled task.
	scheduledHandlerFunc func(context.Context) error

	// scheduledTaskWrapper is a struct that represents a wrapper for a scheduled task.
	// It implements the TaskScheduler and TaskHandler interface.
	scheduledTaskWrapper struct {
		cronSpec string
		name     string
		fn       scheduledHandlerFunc
		opts     []asynq.Option
	}
)

// TaskName returns the name of the task handled by the scheduledTaskWrapper.
func (h *scheduledTaskWrapper) TaskName() string {
	return h.name
}

// Schedule returns the cron specification for the task scheduler.
// It specifies when the task scheduler should be scheduled to run.
func (h *scheduledTaskWrapper) Schedule() string {
	return h.cronSpec
}

// Handle is a method that handles the task scheduler wrapper.
// It takes a context.Context as input and returns an error.
func (h *scheduledTaskWrapper) Handle(ctx context.Context, _ []byte) error {
	return h.fn(ctx)
}

// Options returns the options for the task handler.
func (h *scheduledTaskWrapper) Options() []asynq.Option {
	return h.opts
}

// NewTaskScheduler creates a new task scheduler with the given cron spec and name.
func NewTaskScheduler(cronSpec, name string, opts ...TaskOption) TaskScheduler {
	return &scheduledTaskWrapper{cronSpec: cronSpec, name: name, opts: opts}
}

// ScheduledHandlerFunc is a function that creates a TaskHandler for a scheduled task.
// It takes a name string and a scheduledHandlerFunc as parameters and returns a TaskHandler.
// The name parameter specifies the name of the scheduled task, while the fn parameter is the function to be executed when the task is triggered.
// The returned TaskHandler can be used to register the scheduled task in the queue server.
func ScheduledHandlerFunc(name string, fn scheduledHandlerFunc) TaskHandler {
	return &scheduledTaskWrapper{
		name: name,
		fn:   fn,
	}
}
