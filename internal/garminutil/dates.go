package garminutil

import (
	"fmt"
	"strings"
	"time"
)

// ConfigError marks an invalid plugin configuration value so callers can
// report it to the host as invalid_request instead of provider_unavailable.
type ConfigError struct {
	msg string
}

func (e *ConfigError) Error() string { return e.msg }

var afterLayouts = []string{
	"2006-01-02",
	time.RFC3339,
	"2006-01-02T15:04:05.000Z",
}

// ParseAfterDate reads the optional "after" config option. An absent or empty
// value means "no filter" and returns the zero time with a nil error. A present
// but unparseable value is a ConfigError so the sync fails loudly instead of
// silently importing the full history.
func ParseAfterDate(options map[string]any) (time.Time, error) {
	if options == nil {
		return time.Time{}, nil
	}
	raw, ok := options["after"]
	if !ok || raw == nil {
		return time.Time{}, nil
	}
	afterStr, ok := raw.(string)
	if !ok {
		return time.Time{}, &ConfigError{msg: fmt.Sprintf("invalid after option: expected a date string (YYYY-MM-DD or RFC3339), got %T", raw)}
	}
	afterStr = strings.TrimSpace(afterStr)
	if afterStr == "" {
		return time.Time{}, nil
	}
	for _, layout := range afterLayouts {
		if parsed, err := time.Parse(layout, afterStr); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, &ConfigError{msg: fmt.Sprintf("invalid after option %q: expected YYYY-MM-DD or RFC3339", afterStr)}
}

var activityStartTimeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999-0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
}

// ParseActivityStartTime parses a Garmin activity start time in any of the
// observed response formats.
func ParseActivityStartTime(value string) (time.Time, error) {
	for _, layout := range activityStartTimeLayouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported Garmin activity start time format: %q", value)
}
