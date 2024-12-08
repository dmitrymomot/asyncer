# asyncer

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/asyncer.svg)](https://pkg.go.dev/github.com/dmitrymomot/asyncer)
[![License](https://img.shields.io/github/license/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/blob/main/LICENSE)

[![Tests](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/asyncer)](https://goreportcard.com/report/github.com/dmitrymomot/asyncer)

A type-safe distributed task queue in Go, built on top of [hibiken/asynq](https://github.com/hibiken/asynq).

## Key Features

- Type-safe task handlers
- Support for immediate and scheduled tasks
- Redis-based task queue
- Built-in monitoring UI

## Installation

```bash
go get github.com/dmitrymomot/asyncer
```

## Usage Examples

### Email Service Example

This example demonstrates how to implement an email service with different types of emails:

```go
package email

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/dmitrymomot/asyncer"
    "github.com/redis/go-redis/v9"
)

// Task names
const (
    WelcomeEmailTask     = "email:welcome"
    PasswordResetTask    = "email:password_reset"
    WeeklyDigestTask    = "email:weekly_digest"
)

// Task payloads
type WelcomeEmail struct {
    UserID    int64  `json:"user_id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
}

type PasswordResetEmail struct {
    UserID      int64  `json:"user_id"`
    Email       string `json:"email"`
    ResetToken  string `json:"reset_token"`
    ExpiresAt   int64  `json:"expires_at"`
}

type WeeklyDigestEmail struct {
    UserID       int64    `json:"user_id"`
    Email        string   `json:"email"`
    ArticleIDs   []int64  `json:"article_ids"`
    WeekNumber   int      `json:"week_number"`
}

// EmailService handles email sending operations
type EmailService struct {
    enqueuer *asyncer.Enqueuer
}

// NewEmailService creates a new email service
func NewEmailService(redis *redis.Client) *EmailService {
    return &EmailService{
        enqueuer: asyncer.MustNewEnqueuer(redis),
    }
}

// SendWelcomeEmail enqueues a welcome email task
func (s *EmailService) SendWelcomeEmail(ctx context.Context, userID int64, email, firstName string) error {
    return s.enqueuer.EnqueueTask(ctx, WelcomeEmailTask, WelcomeEmail{
        UserID:    userID,
        Email:     email,
        FirstName: firstName,
    })
}

// SendPasswordResetEmail enqueues a password reset email task
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, userID int64, email, token string, expiresAt int64) error {
    return s.enqueuer.EnqueueTask(ctx, PasswordResetTask, PasswordResetEmail{
        UserID:     userID,
        Email:      email,
        ResetToken: token,
        ExpiresAt:  expiresAt,
    })
}

// ScheduleWeeklyDigest schedules weekly digest emails
func (s *EmailService) ScheduleWeeklyDigest(ctx context.Context, userID int64, email string, articleIDs []int64, weekNum int) error {
    return s.enqueuer.EnqueueTask(ctx, WeeklyDigestTask, WeeklyDigestEmail{
        UserID:     userID,
        Email:      email,
        ArticleIDs: articleIDs,
        WeekNumber: weekNum,
    })
}
```

### Email Worker Example

Implementation of the email processing worker:

```go
package worker

import (
    "context"
    "fmt"

    "github.com/dmitrymomot/asyncer"
    "github.com/redis/go-redis/v9"
    "your/app/email"
    "your/app/mailer" // your email sending implementation
)

type EmailWorker struct {
    mailer mailer.Service
}

func NewEmailWorker(mailer mailer.Service) *EmailWorker {
    return &EmailWorker{mailer: mailer}
}

// HandleWelcomeEmail processes welcome emails
func (w *EmailWorker) HandleWelcomeEmail(ctx context.Context, payload email.WelcomeEmail) error {
    return w.mailer.Send(ctx, mailer.Email{
        To:      payload.Email,
        Subject: "Welcome to Our Platform!",
        Template: "welcome",
        Data: map[string]interface{}{
            "first_name": payload.FirstName,
        },
    })
}

// HandlePasswordReset processes password reset emails
func (w *EmailWorker) HandlePasswordReset(ctx context.Context, payload email.PasswordResetEmail) error {
    return w.mailer.Send(ctx, mailer.Email{
        To:      payload.Email,
        Subject: "Password Reset Request",
        Template: "password_reset",
        Data: map[string]interface{}{
            "reset_link": fmt.Sprintf("https://app.example.com/reset?token=%s", payload.ResetToken),
            "expires_at": payload.ExpiresAt,
        },
    })
}

// HandleWeeklyDigest processes weekly digest emails
func (w *EmailWorker) HandleWeeklyDigest(ctx context.Context, payload email.WeeklyDigestEmail) error {
    articles, err := fetchArticles(ctx, payload.ArticleIDs)
    if err != nil {
        return fmt.Errorf("failed to fetch articles: %w", err)
    }

    return w.mailer.Send(ctx, mailer.Email{
        To:      payload.Email,
        Subject: fmt.Sprintf("Your Weekly Digest - Week %d", payload.WeekNumber),
        Template: "weekly_digest",
        Data: map[string]interface{}{
            "articles": articles,
            "week_number": payload.WeekNumber,
        },
    })
}

// StartWorker initializes and runs the email worker
func StartWorker(ctx context.Context, redis *redis.Client, worker *EmailWorker) error {
    return asyncer.RunQueueServer(
        ctx,
        redis,
        nil, // default logger
        asyncer.HandlerFunc(email.WelcomeEmailTask, worker.HandleWelcomeEmail),
        asyncer.HandlerFunc(email.PasswordResetTask, worker.HandlePasswordReset),
        asyncer.HandlerFunc(email.WeeklyDigestTask, worker.HandleWeeklyDigest),
    )
}
```

### Usage in Application

Example of using the email service in your application:

```go
package main

import (
    "context"

    "your/app/email"
    "your/app/worker"
)

func main() {
    // Initialize your Redis client and other dependencies
    redisClient := initRedis()
    mailerService := initMailer()

    // Initialize email service for enqueueing tasks
    emailService := email.NewEmailService(redisClient)

    // Start the worker in a separate goroutine
    go func() {
        emailWorker := worker.NewEmailWorker(mailerService)
        if err := worker.StartWorker(context.Background(), redisClient, emailWorker); err != nil {
            log.Fatal(err)
        }
    }()

    // Use the email service in your application
    if err := emailService.SendWelcomeEmail(
        context.Background(),
        123,
        "user@example.com",
        "John",
    ); err != nil {
        log.Printf("Failed to send welcome email: %v", err)
    }
}
```

### Scheduled Tasks Example

Example of setting up recurring notifications:

```go
package notifications

import (
    "context"
    "time"

    "github.com/dmitrymomot/asyncer"
    "github.com/redis/go-redis/v9"
)

const DigestSchedulerTask = "scheduler:weekly_digest"

// StartScheduler initializes the task scheduler
func StartScheduler(ctx context.Context, redis *redis.Client) error {
    return asyncer.RunSchedulerServer(
        ctx,
        redis,
        nil, // default logger
        // Schedule weekly digest every Monday at 9 AM
        asyncer.NewTaskScheduler(
            "0 9 * * 1", // cron expression
            DigestSchedulerTask,
            asyncer.Unique(24*time.Hour), // prevent duplicate runs
        ),
    )
}
```

## Configuration

### Enqueuer Options

| Option                              | Description                      |
| ----------------------------------- | -------------------------------- |
| `WithTaskDeadline(d time.Duration)` | Sets maximum task execution time |
| `WithMaxRetry(n int)`               | Sets maximum retry attempts      |
| `WithQueue(name string)`            | Specifies queue name             |
| `WithRetention(d time.Duration)`    | Sets task retention period       |

```go
// Enqueue with options
enqueuer.EnqueueTask(
    ctx,
    taskName,
    payload,
    asyncer.WithMaxRetry(3),
    asyncer.WithQueue("high"),
    asyncer.WithTaskDeadline(5*time.Minute),
)
```

### Scheduler Options

| Option                     | Description                              |
| -------------------------- | ---------------------------------------- |
| `MaxRetry(n int)`          | Sets retry attempts for scheduled tasks  |
| `Timeout(d time.Duration)` | Sets task timeout                        |
| `Unique(d time.Duration)`  | Prevents duplicate tasks within duration |
| `Queue(name string)`       | Specifies queue for scheduled tasks      |

```go
// Schedule with options
asyncer.NewTaskScheduler(
    "*/30 * * * *", // every 30 minutes
    taskName,
    asyncer.MaxRetry(3),
    asyncer.Timeout(1*time.Minute),
    asyncer.Queue("scheduled"),
)
```

## Monitoring

### Asynqmon Web Interface

Access the monitoring dashboard at `http://localhost:8181` to:

- View active, pending, and completed tasks
- Monitor queue statistics
- Inspect task details and errors
- Manage task queues

### Logging

The package supports structured logging through the standard `slog` package:

```go
asyncer.NewSlogAdapter(slog.Default().With(
    slog.String("component", "queue-server"),
))
```

## Contributing

Contributions to the `asyncer` package are welcome! Here are some ways you can contribute:

- Reporting bugs
- **Covering code with tests**
- Suggesting enhancements
- Submitting pull requests
- Sharing the love by telling others about this project

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/dmitrymomot/asyncer/tree/main/LICENSE) file for details. This project contains some code from [hibiken/asynq](https://github.com/hibiken/asynq) package, which is also licensed under the MIT License.
