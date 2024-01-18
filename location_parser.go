package asyncer

import (
	"time"
)

// parseLocation parses the given timeZone string and returns a pointer to a time.Location.
// If the timeZone string is invalid, it returns a pointer to the UTC time.Location.
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
