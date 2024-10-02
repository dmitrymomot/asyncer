package asyncer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
)

type (
	// TaskHandler is an interface for task handlers.
	// It is used to register task handlers in the queue server.
	TaskHandler interface {
		// TaskName returns the name of the task. It is used to register the task handler.
		TaskName() string
		// Handle handles the task. It takes a context and a payload as parameters.
		Handle(ctx context.Context, payload []byte) error
		// Options returns the options for the task handler.
		Options() []asynq.Option
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
		opts     []TaskOption
	}
)

// TaskName returns the name of the task handled by the handlerFuncWrapper.
func (h *handlerFuncWrapper[Payload]) TaskName() string {
	return h.name
}

// Schedule returns the cron specification for the handler function.
// It specifies when the handler function should be scheduled to run.
func (h *handlerFuncWrapper[Payload]) Schedule() string {
	return h.cronSpec
}

// Handle is a method that handles the given payload by unmarshaling it and calling the wrapped handler function.
// It takes a context.Context and a []byte payload as input and returns an error.
// The payload is unmarshaled into a Payload struct, and if the unmarshaling fails, an error is returned.
// Otherwise, the wrapped handler function is called with the context and unmarshaled payload.
func (h *handlerFuncWrapper[Payload]) Handle(ctx context.Context, payload []byte) error {
	var p Payload
	if payload != nil {
		if err := json.Unmarshal(payload, &p); err != nil {
			return errors.Join(ErrFailedToUnmarshalPayload, err)
		}
	}

	return h.fn(ctx, p)
}

// Options returns the options for the handler function.
func (h *handlerFuncWrapper[Payload]) Options() []asynq.Option {
	return h.opts
}

// HandlerFunc is a function that creates a TaskHandler for handling tasks of a specific payload type.
// It takes a name string and a handler function as parameters and returns a TaskHandler.
// The name parameter represents the name of the handler, while the fn parameter is the actual handler function.
// The TaskHandler returned by HandlerFunc is responsible for executing the handler function when a task of the specified payload type is received.
// The payload type is specified using the generic type parameter Payload.
func HandlerFunc[Payload any](name string, fn handlerFunc[Payload], opts ...TaskOption) TaskHandler {
	return &handlerFuncWrapper[Payload]{
		name: name,
		fn:   fn,
		opts: opts,
	}
}
