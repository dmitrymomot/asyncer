package asyncer

import (
	"context"
	"encoding/json"
	"errors"
)

type (
	// TaskHandler is an interface for task handlers.
	// It is used to register task handlers in the queue server.
	TaskHandler interface {
		TaskName() string
		Handle(ctx context.Context, payload []byte) error
	}

	// ScheduledTaskHandler is an interface for scheduled task handlers.
	// It is used to register scheduled task handlers in the scheduler server.
	ScheduledTaskHandler interface {
		TaskHandler
		// Schedule returns the cron spec for the task.
		// For more information about cron spec, see https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format.
		Schedule() string
	}

	// handlerFunc is a function that handles a task.
	handlerFunc[Payload any] func(context.Context, Payload) error

	// handlerFuncWrapper is a struct that represents a wrapper for a handler function.
	// It contains the cronSpec, name, and fn fields.
	// - cronSpec: a string representing the cron specification for scheduling the task.
	// - name: a string representing the name of the task.
	// - fn: a handler function that takes a payload of type Payload and performs the task.
	handlerFuncWrapper[Payload any] struct {
		cronSpec string
		name     string
		fn       handlerFunc[Payload]
	}
)

// TaskName returns the name of the task.
func (h *handlerFuncWrapper[Payload]) TaskName() string {
	return h.name
}

// Schedule returns the cron spec for the task.
func (h *handlerFuncWrapper[Payload]) Schedule() string {
	return h.cronSpec
}

// Handle handles the task.
func (h *handlerFuncWrapper[Payload]) Handle(ctx context.Context, payload []byte) error {
	var p Payload
	if err := json.Unmarshal(payload, &p); err != nil {
		return errors.Join(ErrFailedToUnmarshalPayload, err)
	}

	return h.fn(ctx, p)
}

// HandlerFunc is a function that creates a TaskHandler for handling tasks of a specific payload type.
// It takes a name string and a handler function as parameters and returns a TaskHandler.
// The name parameter represents the name of the handler, while the fn parameter is the actual handler function.
// The TaskHandler returned by HandlerFunc is responsible for executing the handler function when a task of the specified payload type is received.
// The payload type is specified using the generic type parameter Payload.
func HandlerFunc[Payload any](name string, fn handlerFunc[Payload]) TaskHandler {
	return &handlerFuncWrapper[Payload]{
		name: name,
		fn:   fn,
	}
}

// ScheduledHandlerFunc is a function that creates a ScheduledTaskHandler for handling tasks of a specific payload type.
// It takes a name string, a cron spec string, and a handler function as parameters and returns a ScheduledTaskHandler.
// The name parameter represents the name of the handler, while the cronSpec parameter represents the cron spec for the handler.
// The fn parameter is the actual handler function.
// The ScheduledTaskHandler returned by ScheduledHandlerFunc is responsible for executing the handler function when a task of the specified payload type is received.
// The payload type is specified using the generic type parameter Payload.
func ScheduledHandlerFunc[Payload any](cronSpec string, name string, fn handlerFunc[Payload]) ScheduledTaskHandler {
	return &handlerFuncWrapper[Payload]{
		name:     name,
		cronSpec: cronSpec,
		fn:       fn,
	}
}
