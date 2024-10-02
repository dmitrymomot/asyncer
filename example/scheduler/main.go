package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/random"
	"golang.org/x/sync/errgroup"

	"github.com/dmitrymomot/asyncer"
)

const (
	redisAddr    = "redis://localhost:63791/0"
	TestTaskName = "scheduled_task"
)

// test task handler function
func testTaskHandler(handlerN string) func(_ context.Context) error {
	return func(_ context.Context) error {
		// Simulate a task that takes some time to complete.
		fmt.Printf("handler %s: scheduled test task handler called at %s\n", handlerN, time.Now().Format(time.RFC3339))
		time.Sleep(2 * time.Second)
		return nil
	}
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	handlerN := random.String(2, random.Numeric)

	// Run a new queue server with redis as the broker.
	eg.Go(asyncer.RunQueueServer(
		ctx, redisAddr, nil,
		// Register a handler for the task.
		asyncer.ScheduledHandlerFunc(TestTaskName, testTaskHandler(handlerN)),
		// ... add more handlers here ...
	))

	// Run a scheduler with redis as the broker.
	// The scheduler will schedule tasks to be enqueued at a specified time.
	eg.Go(asyncer.RunSchedulerServer(
		ctx, redisAddr, nil,
		// Schedule the scheduled_task task to be enqueued every 1 seconds.
		asyncer.NewTaskScheduler("@every 5s", TestTaskName,
			asyncer.MaxRetry(0),
			asyncer.Timeout(5*time.Second),
			asyncer.Unique(5*time.Second),
		),
		// ... add more scheduled tasks here ...
	))

	// Wait for the queue server to exit.
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
