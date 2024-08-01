package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dmitrymomot/asyncer"
)

const (
	redisAddr     = "redis://localhost:6379/0"
	TestTaskName  = "queued_task"
	TestTaskName2 = "queued_task_2"
)

type TestTaskPayload struct {
	Name string
}

type TestTaskPayload2 struct {
	Greeting string
}

// test task handler function
func testTaskHandler(_ context.Context, payload TestTaskPayload) error {
	fmt.Printf("Hello, %s!\n", payload.Name)
	return nil
}

// test task handler function
func testTaskHandler2(_ context.Context, payload TestTaskPayload2) error {
	fmt.Printf("Hola, %s!\n", payload.Greeting)
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, _ := errgroup.WithContext(ctx)

	// Run a new queue server with redis as the broker.
	eg.Go(asyncer.RunQueueServer(
		ctx, redisAddr,
		asyncer.NewSlogAdapter(slog.Default().With(slog.String("component", "queue-server"))),
		// Register a handler for the task.
		asyncer.HandlerFunc(TestTaskName, testTaskHandler),
		asyncer.HandlerFunc(TestTaskName2, testTaskHandler2),
		// ... add more handlers here ...
	))

	// Create a new enqueuer with redis as the broker.
	enqueuer := asyncer.MustNewEnqueuer(
		redisAddr,
		asyncer.WithTaskDeadline(10*time.Minute),
		asyncer.WithMaxRetry(0),
	)
	defer func(enqueuer *asyncer.Enqueuer) {
		err := enqueuer.Close()
		if err != nil {
			slog.Error("Failed to close the enqueuer", "error", err)
		}
	}(enqueuer)

	// Enqueue a task with payload.
	// The task will be processed after immediately.
	eg.Go(func() error {
		var i int
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				i++
				if err := enqueuer.EnqueueTask(ctx, TestTaskName, TestTaskPayload{
					Name: fmt.Sprintf("Test %d", i),
				}); err != nil {
					return err
				}
				if err := enqueuer.EnqueueTask(ctx, TestTaskName2, TestTaskPayload2{
					Greeting: fmt.Sprintf("Greeter %d", i),
				}); err != nil {
					return err
				}
			}
		}
	})

	// Listen for signals to cancel the context.
	// This will stop the routine and close the queue server.
	eg.Go(func() error {
		c := make(chan os.Signal, 1) // Create channel to signify a signal being sent
		signal.Notify(c, os.Interrupt,
			syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

		select {
		case <-c:
			log.Println("Shutting down...")
			// Cancel the context
			cancel()

			return nil
		case <-ctx.Done():
			return nil
		}
	})

	// Wait for the queue server to exit.
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
