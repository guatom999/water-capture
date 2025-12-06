package utils

import "time"

// FormatTime formats time to ISO 8601
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses ISO 8601 time string
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
