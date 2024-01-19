package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/asyncer"
	"golang.org/x/sync/errgroup"
)

const (
	redisAddr    = "redis://localhost:6379/0"
	TestTaskName = "queued_task"
)

type TestTaskPayload struct {
	Name string
}

// test task handler function
func testTaskHandler(ctx context.Context, payload TestTaskPayload) error {
	fmt.Printf("Hello, %s!\n", payload.Name)
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, _ := errgroup.WithContext(ctx)

	// Run a new queue server with redis as the broker.
	eg.Go(asyncer.RunQueueServer(
		ctx, redisAddr, nil,
		// Register a handler for the task.
		asyncer.HandlerFunc[TestTaskPayload](TestTaskName, testTaskHandler),
		// ... add more handlers here ...
	))

	// Create a new enqueuer with redis as the broker.
	enqueuer := asyncer.MustNewEnqueuer(redisAddr)
	defer enqueuer.Close()

	// Enqueue a task with payload.
	// The task will be processed after immediately.
	for i := 0; i < 10; i++ {
		if err := enqueuer.EnqueueTask(ctx, TestTaskName, TestTaskPayload{
			Name: fmt.Sprintf("Test %d", i),
		}); err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Wait for the queue server to exit.
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
