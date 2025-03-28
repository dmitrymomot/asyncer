package asyncer

import (
	"time"

	"github.com/hibiken/asynq"
)

type TaskOption = asynq.Option

// MaxRetry sets the maximum number of retries for the task.
// The task will be marked as failed after the specified number of failed attempts.
func MaxRetry(n int) TaskOption {
	if n < 0 {
		n = 0
	}
	return asynq.MaxRetry(n)
}

// Timeout sets the timeout for the task.
// The task will be marked as failed if it takes longer than the specified duration.
func Timeout(d time.Duration) TaskOption {
	if d <= 0 {
		d = time.Second
	}
	return asynq.Timeout(d)
}

// Deadline sets the deadline for the task.
// The task will not be processed if it is received after the specified date and time.
func Deadline(t time.Time) TaskOption {
	if t.IsZero() || t.Before(time.Now()) {
		t = time.Now().Add(time.Second)
	}
	return asynq.Deadline(t)
}

// Unique sets the uniqueness constraint for the task.
// The task will not be enqueued if there is an identical task already in the queue.
// The uniqueness constraint is based on the task type and payload.
// The uniqueness constraint is valid for the specified duration.
func Unique(ttl time.Duration) TaskOption {
	if ttl <= 0 {
		ttl = time.Second
	}
	return asynq.Unique(ttl)
}

// TaskID sets the ID for the task.
// The task will be assigned the specified ID.
// Use this option to enqueue a task with a specific ID to prevent duplicate tasks.
// If a task with the same ID already exists in the queue, it will be replaced by the new task.
func TaskID(id string) TaskOption {
	if id != "" {
		return asynq.TaskID(id)
	}
	return nil
}

// Group returns an option to specify the group used for the task.
// Tasks in a given queue with the same group will be aggregated into one task before passed to Handler.
func Group(g string) TaskOption {
	if g != "" {
		return asynq.Group(g)
	}
	return nil
}

// ProcessAt returns an option to specify when to process the given task.
//
// If there's a conflicting ProcessIn option, the last option passed to Enqueue overrides the others.
func ProcessAt(t time.Time) TaskOption {
	if t.IsZero() || t.Before(time.Now()) {
		t = time.Now().Add(time.Second)
	}
	return asynq.ProcessAt(t)
}

// ProcessIn returns an option to specify when to process the given task relative to the current time.
//
// If there's a conflicting ProcessAt option, the last option passed to Enqueue overrides the others.
func ProcessIn(d time.Duration) TaskOption {
	if d <= 0 {
		d = time.Second
	}
	return asynq.ProcessIn(d)
}
