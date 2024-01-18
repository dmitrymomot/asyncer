package asyncer

import "github.com/hibiken/asynq"

// Log levels string representation.
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

// castToAsynqLogLevel converts a string representation of a log level to the corresponding asynq.LogLevel.
// It returns the converted log level or asynq.InfoLevel if the input is not a recognized log level.
func castToAsynqLogLevel(level string) asynq.LogLevel {
	switch level {
	case "debug":
		return asynq.DebugLevel
	case "info":
		return asynq.InfoLevel
	case "warn":
		return asynq.WarnLevel
	case "error":
		return asynq.ErrorLevel
	case "fatal":
		return asynq.FatalLevel
	default:
		return asynq.InfoLevel
	}
}
