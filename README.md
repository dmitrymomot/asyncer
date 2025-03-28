# asyncer

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/asyncer.svg)](https://pkg.go.dev/github.com/dmitrymomot/asyncer)
[![License](https://img.shields.io/github/license/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/blob/main/LICENSE)

[![Tests](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/asyncer)](https://goreportcard.com/report/github.com/dmitrymomot/asyncer)

A type-safe distributed task queue in Go, built on top of [hibiken/asynq](https://github.com/hibiken/asynq) using Redis.

## Key Features

- **Type-safety**: Strongly-typed task handlers and payloads
- **Distributed**: Redis-based task queue for distributed processing
- **Flexible**: Support for immediate and scheduled tasks
- **Configurable**: Extensive options for tasks, queues, and scheduling
- **Monitoring**: Built-in monitoring UI (provided by asynq)
- **Logging**: Integration with Go's standard `slog` package
- **Robust**: Task retries, timeouts, and error handling
- **Performance**: Efficient parallelism based on available CPU cores

## Requirements

- Go 1.23.0 or higher (using toolchain go1.24.1)
- Redis server

## Installation

```bash
go get github.com/dmitrymomot/asyncer
```

## Usage Examples

### Basic Queue Setup

Here's a simple example of setting up a queue server and enqueuing tasks:

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/dmitrymomot/asyncer"
    "github.com/redis/go-redis/v9"
    "golang.org/x/sync/errgroup"
)

// Define task types and payloads
const (
    WelcomeEmailTask = "email:welcome"
)

type WelcomeEmailPayload struct {
    UserID    int64  `json:"user_id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
}

// Define task handler
func welcomeEmailHandler(ctx context.Context, payload WelcomeEmailPayload) error {
    fmt.Printf("Sending welcome email to %s (%s)\n", payload.FirstName, payload.Email)
    // Implement email sending logic here
    return nil
}

func main() {
    // Setup context with cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Create error group
    eg, _ := errgroup.WithContext(ctx)

    // Connect to Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    })
    defer redisClient.Close()

    // Create task enqueuer
    enqueuer := asyncer.MustNewEnqueuer(
        redisClient,
        asyncer.WithQueueNameEnq("default"),
        asyncer.WithTaskDeadline(5 * time.Minute),
        asyncer.WithMaxRetry(3),
    )
    defer enqueuer.Close()

    // Run queue server
    eg.Go(asyncer.RunQueueServer(
        ctx, 
        redisClient,
        asyncer.NewSlogAdapter(slog.Default().With(slog.String("component", "queue-server"))),
        // Register task handlers
        asyncer.HandlerFunc(WelcomeEmailTask, welcomeEmailHandler),
    ))

    // Enqueue a task
    if err := enqueuer.EnqueueTask(
        ctx, 
        WelcomeEmailTask, 
        WelcomeEmailPayload{
            UserID:    123,
            Email:     "user@example.com",
            FirstName: "John",
        },
    ); err != nil {
        fmt.Printf("Failed to enqueue task: %v\n", err)
    }

    // Handle graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    select {
    case <-quit:
        fmt.Println("Shutting down...")
        cancel()
    case <-ctx.Done():
    }

    if err := eg.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

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
    WeeklyDigestTask     = "email:weekly_digest"
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
        enqueuer: asyncer.MustNewEnqueuer(
            redis,
            asyncer.WithQueueNameEnq("default"),
            asyncer.WithTaskDeadline(5 * time.Minute),
            asyncer.WithMaxRetry(3),
        ),
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

### Scheduled Tasks

For scheduled tasks, you can use the scheduler server:

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/dmitrymomot/asyncer"
    "github.com/redis/go-redis/v9"
    "golang.org/x/sync/errgroup"
)

const (
    DailyReportTask = "report:daily"
)

// No payload needed for this scheduled task
func generateDailyReport(ctx context.Context, struct{}) error {
    fmt.Println("Generating daily report...")
    // Implementation of report generation
    return nil
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    eg, _ := errgroup.WithContext(ctx)

    // Connect to Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    })
    defer redisClient.Close()

    // Configure logger
    logger := asyncer.NewSlogAdapter(slog.Default().With(
        slog.String("component", "scheduler-server"),
    ))

    // Run scheduler server - schedules tasks to run
    eg.Go(asyncer.RunSchedulerServer(
        ctx,
        redisClient,
        logger,
        // Schedule daily report at midnight
        asyncer.NewTaskScheduler("0 0 * * *", DailyReportTask),
    ))

    // Run queue server - processes the scheduled tasks
    eg.Go(asyncer.RunQueueServer(
        ctx,
        redisClient,
        logger,
        // Register handler for the scheduled task
        asyncer.HandlerFunc(DailyReportTask, generateDailyReport),
    ))

    // Handle graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    select {
    case <-quit:
        fmt.Println("Shutting down...")
        cancel()
    case <-ctx.Done():
    }

    if err := eg.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Advanced Configuration

### Queue Options

```go
// Configure queue server
queueServer := asyncer.NewQueueServer(
    redisClient,
    // Set worker concurrency
    asyncer.WithQueueConcurrency(10),
    // Set queue priority (higher number = higher priority)
    asyncer.WithQueue("critical", 10),
    asyncer.WithQueue("default", 5),
    asyncer.WithQueue("low", 1),
    // Set worker shutdown timeout
    asyncer.WithQueueShutdownTimeout(30 * time.Second),
    // Set logger
    asyncer.WithQueueLogger(customLogger),
)
```

### Task Options when Initializing Enqueuer

```go
// Configure task options when initializing the Enqueuer
enqueuer := asyncer.MustNewEnqueuer(
    redisClient,
    asyncer.WithQueueNameEnq("default"),
    asyncer.WithTaskDeadline(5 * time.Minute),
    asyncer.WithMaxRetry(3),
)
```

### Task Options when Enqueuing

You can also specify options when enqueuing a task:

```go
// Configure task options when enqueuing
err := enqueuer.EnqueueTask(
    ctx,
    "task:name",
    payload,
    // Set task queue
    asynq.Queue("critical"),
    // Set task processing timeout
    asyncer.Timeout(5 * time.Minute),
    // Schedule task for future execution
    asyncer.ProcessIn(24 * time.Hour),
    // Set retries
    asyncer.MaxRetry(5),
    // Prevent duplicate tasks
    asyncer.Unique(1 * time.Hour),
    // Set task ID
    asyncer.TaskID("unique-task-id"),
    // Set task group
    asyncer.Group("email-notifications"),
    // Set task deadline
    asyncer.Deadline(time.Now().Add(6 * time.Hour)),
)
```

### Scheduler Options

```go
// Configure scheduler server
schedulerServer := asyncer.NewSchedulerServer(
    redisClient,
    // Set timezone
    asyncer.WithSchedulerLocation("UTC"),
    // Set logger
    asyncer.WithSchedulerLogger(customLogger),
)
```

## Logging

The package supports structured logging through the standard `slog` package:

```go
asyncer.NewSlogAdapter(slog.Default().With(
    slog.String("component", "queue-server"),
))
```

## Monitor UI

Asynq provides a web UI for monitoring tasks. You can run it with:

```go
asynqmon.New(asynqmon.Options{
    RedisConnOpt: asynq.RedisClientOpt{Addr: "localhost:6379"},
}).Run(":8080")
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/dmitrymomot/asyncer/tree/main/LICENSE) file for details. This project is built on top of the [hibiken/asynq](https://github.com/hibiken/asynq) package - please refer to their [license](https://github.com/hibiken/asynq/blob/master/LICENSE) for more information.
