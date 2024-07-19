package holos

import (
	"fmt"
	"math"
	"time"
)

// RoundDuration rounds a duration to the nearest unit based on its length.
func RoundDuration(duration time.Duration) string {
	seconds := duration.Seconds()

	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", int(math.Round(seconds)))
	case seconds < 3600:
		minutes := seconds / 60
		return fmt.Sprintf("%dm", int(math.Round(minutes)))
	case seconds < 86400:
		hours := seconds / 3600
		return fmt.Sprintf("%dh", int(math.Round(hours)))
	default:
		days := seconds / 86400
		return fmt.Sprintf("%dd", int(math.Round(days)))
	}
}
