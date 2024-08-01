package asyncer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
)

type slogAdapter struct {
	log *slog.Logger
}

func NewSlogAdapter(log *slog.Logger) asynq.Logger {
	return &slogAdapter{log: log}
}

// Debug logs a message at Debug level.
func (s *slogAdapter) Debug(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	var msg string

	// If there is only one argument, it is the message.
	if len(args) == 1 {
		msg = fmt.Sprint(args[0])
		args = make([]interface{}, 0)
	}

	// If there are more than one argument, the first one is the message.
	if l := len(args); l > 1 && l%2 != 0 {
		msg = fmt.Sprint(args[0])
		args = args[1:]
	}

	s.log.Debug(msg, args...)
}

// Info logs a message at Info level.
func (s *slogAdapter) Info(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	var msg string

	// If there is only one argument, it is the message.
	if len(args) == 1 {
		msg = fmt.Sprint(args[0])
		args = make([]interface{}, 0)
	}

	// If there are more than one argument, the first one is the message.
	if l := len(args); l > 1 && l%2 != 0 {
		msg = fmt.Sprint(args[0])
		args = args[1:]
	}

	s.log.Info(msg, args...)
}

// Warn logs a message at Warning level.
func (s *slogAdapter) Warn(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	var msg string

	// If there is only one argument, it is the message.
	if len(args) == 1 {
		msg = fmt.Sprint(args[0])
		args = make([]interface{}, 0)
	}

	// If there are more than one argument, the first one is the message.
	if l := len(args); l > 1 && l%2 != 0 {
		msg = fmt.Sprint(args[0])
		args = args[1:]
	}

	s.log.Warn(msg, args...)
}

// Error logs a message at Error level.
func (s *slogAdapter) Error(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	var msg string

	// If there is only one argument, it is the message.
	if len(args) == 1 {
		msg = fmt.Sprint(args[0])
		args = make([]interface{}, 0)
	}

	// If there are more than one argument, the first one is the message.
	if l := len(args); l > 1 && l%2 != 0 {
		msg = fmt.Sprint(args[0])
		args = args[1:]
	}

	s.log.Error(msg, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (s *slogAdapter) Fatal(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	var msg string

	// If there is only one argument, it is the message.
	if len(args) == 1 {
		msg = fmt.Sprint(args[0])
		args = make([]interface{}, 0)
	}

	// If there are more than one argument, the first one is the message.
	if l := len(args); l > 1 && l%2 != 0 {
		msg = fmt.Sprint(args[0])
		args = args[1:]
	}

	s.log.Error(msg, args...)
	os.Exit(1)
}
