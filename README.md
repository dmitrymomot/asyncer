# asyncer

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/asyncer.svg)](https://pkg.go.dev/github.com/dmitrymomot/asyncer)
[![License](https://img.shields.io/github/license/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/blob/main/LICENSE)


[![Tests](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/asyncer)](https://goreportcard.com/report/github.com/dmitrymomot/asyncer)

This is a simple, reliable, and efficient distributed task queue in Go.
The asyncer just wrapps [hibiken/asynq](https://github.com/hibiken/asynq) package with some predefined settings. So, if you need more flexibility, you can use [hibiken/asynq](https://github.com/hibiken/asynq) directly.

## Installation

To install the `asyncer` package, use the following command:

```bash
go get github.com/dmitrymomot/asyncer
```

## Usage

In this example, we will create a simple task that prints a greeting message to the console.
Also, we will create a scheduler that will run the task every hour.

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/dmitrymomot/asyncer"
)

const redisAddr = "redis://localhost:6379/0"
const TestTaskName = "testTask1"

type TestTaskPayload struct {
    Name string
}

// test task handler function
func testTaskHandler(ctx context.Context, payload TestTaskPayload) error {
    fmt.Printf("Hello, %s!\n", payload.Name)
    return nil
}

func main() {
    // Create a new enqueuer with redis as the broker.
    enqueuer := asyncer.MustNewEnqueuer(redisAddr)
    defer enqueuer.Close()

    // Enqueue a task with payload.
    // The task will be processed after immediately.
    if err := enqueuer.EnqueueTask(context.TODO(), TestTaskName, testTaskPayload{Name: "test"}); err != nil {
        panic(err)
    }

    eg, _ := errgroup.WithContext(context.Background())

    // Run a new queue server with redis as the broker.
    eg.Go(asyncer.RunQueueServer(
        redisAddr, nil, 
        // Register a handler for the task.
        asyncer.HandlerFunc[TestTaskPayload](TestTaskName, testTaskHandler),
        // ... add more handlers here ...
    ))

    // Run a scheduler with redis as the broker.
    // The scheduler will schedule tasks to be enqueued at a specified time. 
	eg.Go(asyncer.RunSchedulerServer(
		redisAddr, nil,
        // Schedule the testTask1 task to be enqueued every hour.
		asyncer.ScheduledHandlerFunc[Payload]("@every 1h", TestTaskName),
        // ... add more scheduled tasks here ...
	))
    
    // Wait for the queue server to exit.
    if err := eg.Wait(); err != nil {
        panic(err)
    }
}
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