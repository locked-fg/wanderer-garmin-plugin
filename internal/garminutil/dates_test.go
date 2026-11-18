package garminutil

import (
	"errors"
	"testing"
	"time"
)

func TestParseAfterDateAbsentMeansNoFilter(t *testing.T) {
	cases := map[string]map[string]any{
		"nil options":  nil,
		"key missing":  {},
		"null value":   {"after": nil},
		"empty string": {"after": ""},
		"blank string": {"after": "   "},
	}
	for name, options := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := ParseAfterDate(options)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.IsZero() {
				t.Fatalf("expected zero time, got %v", got)
			}
		})
	}
}

func TestParseAfterDateValidFormats(t *testing.T) {
	cases := map[string]time.Time{
		"2024-01-02":                time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		"2024-01-02T15:04:05Z":      time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		"2024-01-02T15:04:05.000Z":  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		"2024-01-02T15:04:05+02:00": time.Date(2024, 1, 2, 13, 4, 5, 0, time.UTC),
		" 2024-01-02 ":              time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			got, err := ParseAfterDate(map[string]any{"after": input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func TestParseAfterDateInvalidValues(t *testing.T) {
	cases := map[string]any{
		"2024-13-01":          "2024-13-01",
		"german date":         "01.02.2024",
		"datetime with space": "2024-01-02 10:00:00",
		"not a date":          "yesterday",
		"numeric value":       42.0,
		"boolean value":       true,
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := ParseAfterDate(map[string]any{"after": value})
			if err == nil {
				t.Fatalf("expected error, got %v", got)
			}
			var configErr *ConfigError
			if !errors.As(err, &configErr) {
				t.Fatalf("expected *ConfigError, got %T: %v", err, err)
			}
		})
	}
}

func TestParseActivityStartTimeValidFormats(t *testing.T) {
	values := []string{
		"2024-01-02T15:04:05.123456789Z",
		"2024-01-02T15:04:05.123456789+0200",
		"2024-01-02T15:04:05+0200",
		"2024-01-02T15:04:05.5",
		"2024-01-02T15:04:05.000Z",
		"2024-01-02T15:04:05Z",
		"2024-01-02T15:04:05.000",
		"2024-01-02T15:04:05",
		"2024-01-02 15:04:05.25",
		"2024-01-02 15:04:05",
	}
	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			got, err := ParseActivityStartTime(value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Year() != 2024 || got.Month() != 1 || got.Day() != 2 {
				t.Fatalf("unexpected date: %v", got)
			}
		})
	}
}

func TestParseActivityStartTimeRejectsInvalid(t *testing.T) {
	for _, value := range []string{"", "not a date", "2024/01/02", "2024-01-02T15:04:05+02:00:00"} {
		t.Run(value, func(t *testing.T) {
			if _, err := ParseActivityStartTime(value); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
