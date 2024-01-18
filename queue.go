package asyncer

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/hibiken/asynq"
)

type (
	// QueueServer is a wrapper for asynq.Server.
	QueueServer struct {
		asynq *asynq.Server
	}

	// QueueServerOption is a function that configures a QueueServer.
	QueueServerOption func(*asynq.Config)

	// TaskHandler is an interface for task handlers.
	// It is used to register task handlers in the queue server.
	TaskHandler interface {
		TaskName() string
		Handle(ctx context.Context, payload []byte) error
	}
)

// NewQueueServer creates a new instance of QueueServer.
// It takes a redis connection option and optional queue server options.
// The function returns a pointer to the created QueueServer.
func NewQueueServer(redisConnOpt asynq.RedisConnOpt, opts ...QueueServerOption) *QueueServer {
	// Get the number of available CPUs.
	useProcs := runtime.GOMAXPROCS(0)
	if useProcs == 0 {
		useProcs = 1
	} else if useProcs > 1 {
		useProcs = useProcs / 2
	}

	// Default queue options
	var (
		workerConcurrency     = useProcs // use half of the available CPUs
		workerShutdownTimeout = time.Second * 10
		workerLogLevel        = "info"
		queueName             = "default"
	)

	cnf := asynq.Config{
		Concurrency:     workerConcurrency,
		LogLevel:        castToAsynqLogLevel(workerLogLevel),
		ShutdownTimeout: workerShutdownTimeout,
		Queues: map[string]int{
			queueName: workerConcurrency,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(&cnf)
	}

	return &QueueServer{asynq: asynq.NewServer(redisConnOpt, cnf)}
}

// Run starts the queue server and registers the provided task handlers.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(queueServer.Run(
//		NewTaskHandler1(),
//		NewTaskHandler2(),
//	))
//
// The function returns an error if the server fails to start.
func (srv *QueueServer) Run(handlers ...TaskHandler) func() error {
	return func() error {
		mux := asynq.NewServeMux()

		// Register handlers
		for _, h := range handlers {
			mux.HandleFunc(h.TaskName(), func(ctx context.Context, t *asynq.Task) error {
				return h.Handle(ctx, t.Payload())
			})
		}

		// Run server
		if err := srv.asynq.Run(mux); err != nil {
			return errors.Join(ErrFailedToStartQueueServer, err)
		}

		return nil
	}
}

// Shutdown gracefully shuts down the queue server by waiting for all
// in-flight tasks to finish processing before shutdown.
func (srv *QueueServer) Shutdown() {
	srv.asynq.Stop()
	srv.asynq.Shutdown()
}
