package asyncer

import "github.com/hibiken/asynq"

// WithSchedulerLogLevel sets the scheduler log level.
func WithSchedulerLogLevel(level string) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		cnf.LogLevel = getAsynqLogLevel(level)
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
		cnf.PreEnqueueFunc = fn
	}
}

// WithPostEnqueueFunc sets the scheduler post enqueue function.
func WithPostEnqueueFunc(fn func(info *asynq.TaskInfo, err error)) SchedulerServerOption {
	return func(cnf *asynq.SchedulerOpts) {
		cnf.PostEnqueueFunc = fn
	}
}
