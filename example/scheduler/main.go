package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dmitrymomot/asyncer"
)

const (
	redisAddr    = "redis://localhost:6379/0"
	TestTaskName = "scheduled_task"
)

// test task handler function
func testTaskHandler(_ context.Context) error {
	fmt.Println("scheduled test task handler called at", time.Now().Format(time.RFC3339))
	return nil
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())

	// Run a new queue server with redis as the broker.
	eg.Go(asyncer.RunQueueServer(
		ctx, redisAddr, nil,
		// Register a handler for the task.
		asyncer.ScheduledHandlerFunc(TestTaskName, testTaskHandler),
		// ... add more handlers here ...
	))

	// Run a scheduler with redis as the broker.
	// The scheduler will schedule tasks to be enqueued at a specified time.
	eg.Go(asyncer.RunSchedulerServer(
		ctx, redisAddr, nil,
		// Schedule the scheduled_task task to be enqueued every 1 seconds.
		asyncer.NewTaskScheduler("@every 1s", TestTaskName),
		// ... add more scheduled tasks here ...
	))

	// Wait for the queue server to exit.
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
