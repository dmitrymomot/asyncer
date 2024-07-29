package asyncer

import (
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
)

type slogAdapter struct {
	log slog.Logger
}

func NewSlogAdapter(log slog.Logger) asynq.Logger {
	return &slogAdapter{log: log}
}

// Debug logs a message at Debug level.
func (s *slogAdapter) Debug(args ...interface{}) {
	s.log.Debug("", args...)
}

// Info logs a message at Info level.
func (s *slogAdapter) Info(args ...interface{}) {
	s.log.Info("", args...)
}

// Warn logs a message at Warning level.
func (s *slogAdapter) Warn(args ...interface{}) {
	s.log.Warn("", args...)
}

// Error logs a message at Error level.
func (s *slogAdapter) Error(args ...interface{}) {
	s.log.Error("", args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (s *slogAdapter) Fatal(args ...interface{}) {
	s.log.Error("", args...)
	os.Exit(1)
}
