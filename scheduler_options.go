package asyncer

import "github.com/hibiken/asynq"

// WithSchedulerLogLevel sets the scheduler log level.
func WithSchedulerLogLevel(level string) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		cnf.LogLevel = castToAsynqLogLevel(level)
	}
}

// WithSchedulerLogger sets the scheduler logger.
func WithSchedulerLogger(logger asynq.Logger) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		if logger != nil {
			cnf.Logger = logger
		}
	}
}

// WithSchedulerLocation sets the scheduler location.
func WithSchedulerLocation(timeZone string) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		// parse location from string and set it to the config
		cnf.Location = parseLocation(timeZone)
	}
}

// WithPreEnqueueFunc sets the scheduler pre enqueue function.
func WithPreEnqueueFunc(fn func(task *asynq.Task, opts []asynq.Option)) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		if fn != nil {
			cnf.PreEnqueueFunc = fn
		}
	}
}

// WithPostEnqueueFunc sets the scheduler post enqueue function.
func WithPostEnqueueFunc(fn func(info *asynq.TaskInfo, err error)) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		if fn != nil {
			cnf.PostEnqueueFunc = fn
		}
	}
}
