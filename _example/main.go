package main

import (
	"github.com/dmitrymomot/asyncer"
	utils "github.com/dmitrymomot/go-utils"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Init logger
	logger := logrus.WithFields(logrus.Fields{
		"app":       "asyncer",
		"component": "main",
	})
	defer func() { logger.Info("Server successfully shutdown") }()

	// Redis connect options for asynq client
	redisConnOpt, err := asynq.ParseRedisURI("redis://localhost:6379/0") // TODO: Add your redis connection string
	if err != nil {
		logger.WithError(err).Fatal("failed to parse redis connection string")
	}

	// Init asynq client
	asynqClient := asynq.NewClient(redisConnOpt)
	defer asynqClient.Close()

	// Create a context with a timeout and set the Server's context
	ctx, cancel := utils.NewContextWithCancel(logger.WithField("component", "context"))
	defer cancel()

	// Create a new errgroup
	eg, _ := errgroup.WithContext(ctx)

	// Create a new scheduler server with the given options
	schedulerServer := asyncer.NewSchedulerServer(
		redisConnOpt, logger,
		asyncer.WithSchedulerLocation("UTC"), // options are not required
	)
	defer schedulerServer.Shutdown()

	// Init scheduler handlers
	testScheduler := NewScheduler(nil)
	// Run the scheduler
	eg.Go(schedulerServer.Run(
		testScheduler,
		// TODO: Add more schedulers here
	))

	// Create a new queue worker server with the given options
	queueServer := asyncer.NewQueueServer(
		redisConnOpt, logger,
		asyncer.WithQueueName("default"), // options are not required
	)
	defer queueServer.Shutdown()

	// Init worker handlers
	testWorker := NewWorker(nil)
	// Run the worker
	eg.Go(queueServer.Run(
		testWorker,
		// TODO: Add more workers here
	))

	// Create a new enqueuer
	enqueuer := NewEnqueuer(asyncer.NewEnqueuer(asynqClient))
	// TODO: Use enqueuer to enqueue tasks in your app,
	// E,g: enqueuer.SendConfirmationEmail(ctx, "test@example", "123456")
	_ = enqueuer

	// Wait for the server to finish
	if err := eg.Wait(); err != nil {
		logger.WithError(err).Error("Server stopped with error")
	}
}
