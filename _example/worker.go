package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type (
	// Worker is a task handler for email delivery.
	Worker struct {
		mail mailer
	}

	mailer interface {
		SendConfirmationEmail(ctx context.Context, email, otp string) error
	}
)

// NewWorker creates a new email task handler.
func NewWorker(mail mailer) *Worker {
	return &Worker{mail: mail}
}

// Register registers task handlers for email delivery.
func (w *Worker) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(SendConfirmationEmailTask, w.SendConfirmationEmail)
}

// SendConfirmationEmail sends an example confirmation email.
// It is a task handler for SendConfirmationEmailTask.
func (w *Worker) SendConfirmationEmail(ctx context.Context, t *asynq.Task) error {
	var p SendConfirmationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := w.mail.SendConfirmationEmail(ctx, p.Email, p.OTP); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
