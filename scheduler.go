package asyncer

import (
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
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
func NewSchedulerServer(redisClient redis.UniversalClient, opts ...SchedulerServerOption) *SchedulerServer {
	// setup asynq scheduler config
	cnf := &asynq.SchedulerOpts{
		LogLevel: asynq.ErrorLevel,
		Location: time.UTC,
	}

	// Apply options
	for _, opt := range opts {
		opt(cnf)
	}

	return &SchedulerServer{
		asynq: asynq.NewSchedulerFromRedisClient(redisClient, cnf),
	}
}

// ScheduleTask schedules a task based on the given cron specification and task name.
// It returns an error if the cron specification or task name is empty, or if there was an error registering the task.
func (srv *SchedulerServer) ScheduleTask(cronSpec, taskName string, opts ...TaskOption) error {
	if cronSpec == "" {
		return errors.Join(ErrFailedToScheduleTask, ErrCronSpecIsEmpty)
	}
	if taskName == "" {
		return errors.Join(ErrFailedToScheduleTask, ErrTaskNameIsEmpty)
	}

	if _, err := srv.asynq.Register(cronSpec, asynq.NewTask(taskName, nil, opts...)); err != nil {
		return errors.Join(ErrFailedToScheduleTask, err)
	}

	return nil
}

// Run runs the scheduler with the provided handlers.
// It returns a function that can be used to run server in a error group.
// E.g.:
//
//	eg, ctx := errgroup.WithContext(context.Background())
//	eg.Go(schedulerServer.Run())
func (srv *SchedulerServer) Run() func() error {
	return func() error {
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
//		asyncer.NewTaskScheduler("@every 1h", "scheduled_task_1"),
//	))

//	eg.Go(asyncer.RunQueueServer(
//		"redis://localhost:6379",
//		logger,
//		asyncer.ScheduledHandlerFunc("scheduled_task_1", scheduledTaskHandler),
//	))
//
//	func scheduledTaskHandler(ctx context.Context) error {
//		// ...handle task here...
//	}
//
// The function returns an error if the server fails to start.
// The function panics if the Redis connection string is invalid.
//
// !!! Pay attention, that the scheduler just triggers the job, so you need to run queue server as well.
func RunSchedulerServer(ctx context.Context, redisClient redis.UniversalClient, log asynq.Logger, schedulers ...TaskScheduler) func() error {
	// Init scheduler server
	var opts []SchedulerServerOption
	if log != nil {
		opts = append(opts, WithSchedulerLogger(log))
	}

	return func() error {
		srv := NewSchedulerServer(redisClient, opts...)
		defer srv.Shutdown()

		// Register schedulers
		for _, scheduler := range schedulers {
			if err := srv.ScheduleTask(scheduler.Schedule(), scheduler.TaskName(), scheduler.Options()...); err != nil {
				return errors.Join(ErrFailedToRunSchedulerServer, err)
			}
		}

		// Run server
		eg, _ := errgroup.WithContext(ctx)
		eg.Go(srv.Run())
		return eg.Wait()
	}
}
