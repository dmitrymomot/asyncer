package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/asyncer"
	"github.com/dmitrymomot/random"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

const (
	redisAddr    = "localhost:63791"
	TestTaskName = "scheduled_task"
)

// test task handler function
func testTaskHandler(handlerN string) func(_ context.Context) error {
	return func(_ context.Context) error {
		fmt.Printf("handler %s: scheduled test task handler called at %s\n", handlerN, time.Now().Format(time.RFC3339))
		time.Sleep(2 * time.Second)
		return nil
	}
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	handlerN := random.String(2, random.Numeric)

	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
	defer redisClient.Close()

	// Run a new queue server with redis as the broker.
	eg.Go(asyncer.RunQueueServer(
		ctx, redisClient, nil,
		// Register a handler for the task.
		asyncer.ScheduledHandlerFunc(TestTaskName, testTaskHandler(handlerN)),
		// ... add more handlers here ...
	))

	// Run a scheduler with redis as the broker.
	eg.Go(asyncer.RunSchedulerServer(
		ctx, redisClient, nil,
		// Schedule the scheduled_task task to be enqueued every 5 seconds.
		asyncer.NewTaskScheduler("@every 5s", TestTaskName,
			asyncer.MaxRetry(0),
			asyncer.Timeout(5*time.Second),
			asyncer.Unique(5*time.Second),
		),
		// ... add more scheduled tasks here ...
	))

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
