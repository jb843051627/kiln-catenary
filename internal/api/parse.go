package api

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func Time(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("time is required")
	}
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time: %w", err)
	}
	return t.UTC(), nil
}
func Bool(values url.Values, name string, fallback bool) bool {
	value := values.Get(name)
	if value == "" {
		return fallback
	}
	return value == "1" || value == "true"
}
func Limit(values url.Values, fallback int) (int, error) {
	value := values.Get("limit")
	if value == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 || n > 1000 {
		return 0, fmt.Errorf("limit must be between 1 and 1000")
	}
	return n, nil
}
