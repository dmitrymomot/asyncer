package asyncer

import (
	"runtime"
	"time"

	"github.com/hibiken/asynq"
)

type (
	// QueueServer is a wrapper for asynq.Server.
	QueueServer struct {
		*asynq.Server
	}

	// QueueServerOption is a function that configures a QueueServer.
	QueueServerOption func(*asynq.Config)

	// taskHandler is an interface for task handlers.
	taskHandler interface {
		Register(*asynq.ServeMux)
	}
)

// NewQueueServer creates a new queue client and returns the server.
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
		LogLevel:        getAsynqLogLevel(workerLogLevel),
		ShutdownTimeout: workerShutdownTimeout,
		Queues: map[string]int{
			queueName: workerConcurrency,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(&cnf)
	}

	return &QueueServer{Server: asynq.NewServer(redisConnOpt, cnf)}
}

// Run  creates a new queue client, registers task handlers and runs the server.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(queueServer.Run(
//		NewTaskHandler1(),
//		NewTaskHandler2(),
//	))
func (srv *QueueServer) Run(handlers ...taskHandler) func() error {
	return func() error {
		// Run server
		return srv.Server.Run(registerQueueHandlers(handlers...))
	}
}

// registerQueueHandlers registers handlers for each task type.
func registerQueueHandlers(handlers ...taskHandler) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	// Register handlers
	for _, h := range handlers {
		h.Register(mux)
	}

	return mux
}

// Shutdown gracefully shuts down the queue server by waiting for all
// in-flight tasks to finish processing before shutdown.
func (srv *QueueServer) Shutdown() {
	srv.Server.Stop()
	srv.Server.Shutdown()
}
