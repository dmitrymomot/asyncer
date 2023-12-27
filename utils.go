package asyncer

import (
	"time"

	"github.com/hibiken/asynq"
)

// get asynq log level by string.
func getAsynqLogLevel(level string) asynq.LogLevel {
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

// parseLocation parses a location from a string.
//
// If the name is "" or "UTC", LoadLocation returns UTC.
// If the name is "Local", LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
func parseLocation(timeZone string) *time.Location {
	// parse location from string and set it to the config
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.UTC
	}
	return loc
}
