package asyncer

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type (
	// QueueServer is a wrapper for asynq.Server.
	QueueServer struct {
		asynq *asynq.Server
	}

	// QueueServerOption is a function that configures a QueueServer.
	QueueServerOption func(*asynq.Config)
)

// NewQueueServer creates a new instance of QueueServer.
// It takes a redis connection option and optional queue server options.
// The function returns a pointer to the created QueueServer.
func NewQueueServer(redisClient redis.UniversalClient, opts ...QueueServerOption) *QueueServer {
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

	return &QueueServer{asynq: asynq.NewServerFromRedisClient(redisClient, cnf)}
}

// Run starts the queue server and registers the provided task handlers.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(queueServer.Run(
//		yourapp.NewTaskHandler1(),
//		yourapp.NewTaskHandler2(),
//	))
//
// The function returns an error if the server fails to start.
func (srv *QueueServer) Run(handlers ...TaskHandler) func() error {
	return func() error {
		mux := asynq.NewServeMux()

		// Register handlers
		for _, h := range handlers {
			handlerFunc := func(
				fn func(ctx context.Context, payload []byte) error,
			) func(ctx context.Context, t *asynq.Task) error {
				return func(ctx context.Context, t *asynq.Task) error {
					return fn(ctx, t.Payload())
				}
			}
			mux.HandleFunc(h.TaskName(), handlerFunc(h.Handle))
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

// RunQueueServer starts the queue server and registers the provided task handlers.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, _ := errgroup.WithContext(context.Background())
//	eg.Go(asyncer.RunQueueServer(
//		"redis://localhost:6379",
//		logger,
//		asyncer.HandlerFunc[PayloadStruct1]("task1", task1Handler),
//		asyncer.HandlerFunc[PayloadStruct2]("task2", task2Handler),
//	))
//
//	func task1Handler(ctx context.Context, payload PayloadStruct1) error {
//		// ... handle task here ...
//	}
//
//	func task2Handler(ctx context.Context, payload PayloadStruct2) error {
//		// ... handle task here ...
//	}
//
// The function panics if the redis connection string is invalid.
// The function returns an error if the server fails to start.
func RunQueueServer(ctx context.Context, redisClient redis.UniversalClient, log asynq.Logger, handlers ...TaskHandler) func() error {
	// Queue server options
	var opts []QueueServerOption
	if log != nil {
		opts = append(opts, WithQueueLogger(log))
	}

	return func() error {
		srv := NewQueueServer(redisClient, opts...)
		defer srv.Shutdown()

		// Run server
		eg, _ := errgroup.WithContext(ctx)
		eg.Go(srv.Run(handlers...))
		return eg.Wait()
	}
}
