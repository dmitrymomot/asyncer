package asyncer

import "errors"

// Predefined errors.
var (
	ErrFailedToParseRedisURI            = errors.New("failed to parse redis connection string")
	ErrMissedAsynqClient                = errors.New("missed asynq client")
	ErrFailedToCreateEnqueuerWithClient = errors.New("failed to create enqueuer with asynq client")
	ErrFailedToEnqueueTask              = errors.New("failed to enqueue task")
	ErrFailedToCloseEnqueuer            = errors.New("failed to close enqueuer")
	ErrFailedToStartQueueServer         = errors.New("failed to start queue server")
	ErrFailedToUnmarshalPayload         = errors.New("failed to unmarshal payload")
	ErrFailedToRunQueueServer           = errors.New("failed to run queue server")
)
