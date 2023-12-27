package asyncer

import (
	"time"

	"github.com/hibiken/asynq"
)

type (
	// SchedulerServer is a wrapper for asynq.Scheduler.
	SchedulerServer struct {
		*asynq.Scheduler
	}

	// SchedulerServerOption is a function that configures a SchedulerServer.
	SchedulerServerOption func(*asynq.SchedulerOpts)

	// schedulerHandler is an interface for scheduler handlers.
	schedulerHandler interface {
		Schedule(*asynq.Scheduler) error
	}
)

// NewSchedulerServer creates a new scheduler client and returns the server.
func NewSchedulerServer(redisConnOpt asynq.RedisConnOpt, log asynq.Logger, opts ...SchedulerServerOption) *SchedulerServer {
	// Default scheduler options
	var (
		workerLogLevel = "info"
		timeZone       = time.UTC
	)

	// setup asynq scheduler config
	cnf := &asynq.SchedulerOpts{
		Logger:   log,
		LogLevel: getAsynqLogLevel(workerLogLevel),
		Location: timeZone,
	}

	// Apply options
	for _, opt := range opts {
		opt(cnf)
	}

	return &SchedulerServer{Scheduler: asynq.NewScheduler(redisConnOpt, cnf)}
}

// Run scheduler server.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(schedulerServer.Run(
//		NewSchedulerHandler1(),
//		NewSchedulerHandler2(),
//	))
func (srv *SchedulerServer) Run(handlers ...schedulerHandler) func() error {
	return func() error {
		// Register handlers
		for _, h := range handlers {
			if err := h.Schedule(srv.Scheduler); err != nil {
				return err
			}
		}

		// Run scheduler
		return srv.Scheduler.Run()
	}
}

// Shutdown gracefully shuts down the scheduler server by waiting for all
// pending tasks to be processed.
func (srv *SchedulerServer) Shutdown() {
	srv.Scheduler.Shutdown()
}
