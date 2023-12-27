package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
)

// dummy error that means there are no expired OTPs.
// it's here just to show how to avoid retrying the task.
var ErrNoExpiredOTP = errors.New("no expired OTPs")

type (
	// Scheduled tasks runner.
	Scheduler struct {
		repo repository
	}

	// dummy repository interface.
	repository interface {
		CleanUpExpiredOTP() error
	}
)

// NewScheduler creates a new scheduler.
func NewScheduler(repo repository) *Scheduler {
	return &Scheduler{repo: repo}
}

// Schedule schedules tasks for the asyncer.
func (sch *Scheduler) Schedule(s *asynq.Scheduler) error {
	if _, err := s.Register("@every 1h", asynq.NewTask(CleanUpExpiredOTP, nil)); err != nil {
		return fmt.Errorf("failed to register task %s: %w", CleanUpExpiredOTP, err)
	}
	return nil
}

// CleanUpExpiredOTP cleans up expired OTPs.
// It is a task handler for CleanUpExpiredOTP.
func (s *Scheduler) CleanUpExpiredOTP(ctx context.Context, t *asynq.Task) error {
	if err := s.repo.CleanUpExpiredOTP(); err != nil {
		// This way we can avoid retrying the task if there are no expired OTPs.
		if err != ErrNoExpiredOTP {
			return fmt.Errorf("failed to clean up expired OTPs: %w", err)
		}
	}

	return nil
}
