package asyncer

import (
	"errors"
	"time"

	"github.com/hibiken/asynq"
)

type (
	// SchedulerServer is a wrapper for asynq.Scheduler.
	SchedulerServer struct {
		asynq *asynq.Scheduler
	}

	// SchedulerServerOption is a function that configures a SchedulerServer.
	SchedulerServerOption func(*asynq.SchedulerOpts)
)

// NewSchedulerServer creates a new scheduler client and returns the server.
func NewSchedulerServer(redisConnOpt asynq.RedisConnOpt, opts ...SchedulerServerOption) *SchedulerServer {
	// setup asynq scheduler config
	cnf := &asynq.SchedulerOpts{
		LogLevel: asynq.ErrorLevel,
		Location: time.UTC,
	}

	// Apply options
	for _, opt := range opts {
		opt(cnf)
	}

	return &SchedulerServer{asynq: asynq.NewScheduler(redisConnOpt, cnf)}
}

// Run runs the scheduler with the provided handlers.
// It registers each handler with the scheduler and then starts the scheduler.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(schedulerServer.Run(
//		NewSchedulerHandler1(),
//		NewSchedulerHandler2(),
//	))
func (srv *SchedulerServer) Run(handlers ...ScheduledTaskHandler) func() error {
	return func() error {
		// Register handlers
		for _, h := range handlers {
			_, err := srv.asynq.Register(h.Schedule(), asynq.NewTask(h.TaskName(), nil))
			if err != nil {
				return errors.Join(ErrFailedToScheduleTask, err)
			}
		}

		// Run scheduler
		if err := srv.asynq.Run(); err != nil {
			return errors.Join(ErrFailedToStartSchedulerServer, err)
		}
		return nil
	}
}

// Shutdown gracefully shuts down the scheduler server by waiting for all
// pending tasks to be processed.
func (srv *SchedulerServer) Shutdown() {
	srv.asynq.Shutdown()
}

// RunSchedulerServer runs the scheduler server with the given Redis connection string,
// logger, and scheduled task handlers.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(asyncer.RunSchedulerServer(
//		"redis://localhost:6379",
//		logger,
//		asyncer.ScheduledHandlerFunc[Payload]("@every 1h", "scheduled_task_1"),
//	))

//	eg.Go(asyncer.RunQueueServer(
//		"redis://localhost:6379",
//		logger,
//		asyncer.HandlerFunc[Payload]("scheduled_task_1", scheduledTaskHandler),
//	))
//
//	func scheduledTaskHandler(ctx context.Context, payload Payload) error {
//		// ...handle task here...
//	}
//
// The function returns an error if the server fails to start.
// The function panics if the Redis connection string is invalid.
//
// !!! Pay attention, that the scheduler just triggers the job, so you need to run queue server as well.
func RunSchedulerServer(redisConnStr string, log asynq.Logger, handlers ...ScheduledTaskHandler) func() error {
	// Redis connect options for asynq client
	redisConnOpt, err := asynq.ParseRedisURI(redisConnStr)
	if err != nil {
		panic(errors.Join(ErrFailedToRunQueueServer, err))
	}

	// Init scheduler server
	opts := []SchedulerServerOption{}
	if log != nil {
		opts = append(opts, WithSchedulerLogger(log))
	}

	// Init scheduler server
	return NewSchedulerServer(redisConnOpt, opts...).Run(handlers...)
}
