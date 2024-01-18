package asyncer

import (
	"context"
	"encoding/json"
	"errors"
)

// handlerFunc is a function that handles a task.
type handlerFunc[Payload any] func(context.Context, Payload) error

// handlerFuncWrapper is a wrapper for handlerFunc.
// It implements taskHandler interface.
type handlerFuncWrapper[Payload any] struct {
	name string
	fn   handlerFunc[Payload]
}

// TaskName returns the name of the task.
func (h *handlerFuncWrapper[Payload]) TaskName() string {
	return h.name
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
