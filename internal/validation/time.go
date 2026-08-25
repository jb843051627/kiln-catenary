package validation

import (
	"fmt"
	"time"
)

func Ordered(start, end time.Time) error {
	if start.IsZero() || end.IsZero() || !end.After(start) {
		return fmt.Errorf("time range is invalid")
	}
	return nil
}

func InStage(at, start, end time.Time) bool { return !at.Before(start) && at.Before(end) }

func Fresh(at, now time.Time, maxAge time.Duration) bool {
	return !at.After(now) && now.Sub(at) <= maxAge
}
