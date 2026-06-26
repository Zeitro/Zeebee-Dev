package timex

import (
	"fmt"
	"strings"
	"time"
)

// StringToDuration parse string to duration

func StringToDuration(str string) (t time.Duration, err error) {
	return time.ParseDuration(strings.TrimSpace(str))
}

// DurationToString parse duration to strings.
// When humanReadable is false, returns the short form (e.g. "5m", "1h30m45s").
// When humanReadable is true, returns the human-readable form (e.g. "1 hour 30 minutes 45 seconds").

func DurationToString(d time.Duration, humanReadable bool) string {
	if humanReadable {
		return longDuration(d)
	}
	return shortDuration(d)
}

// shortDuration returns the compact Go-style form with zero components
// omitted: "10m0s" -> "10m", "2h0m30s" -> "2h30s", "1h30m45s" unchanged.
func shortDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}
	// Sub-second precision (ms/us/ns) can't be represented in the h/m/s form
	// without dropping the fraction; fall back to stdlib formatting instead.
	if d%time.Second != 0 {
		return d.String()
	}
	neg := d < 0
	if neg {
		d = -d
	}
	hours := d / time.Hour
	d %= time.Hour
	minutes := d / time.Minute
	d %= time.Minute
	seconds := d / time.Second

	// Emit only non-zero units largest-to-smallest. d is non-zero and
	// whole-second here, so at least one unit is appended.
	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	s := strings.Join(parts, "")
	if neg {
		return "-" + s
	}
	return s
}

func longDuration(d time.Duration) string {
	if d <= 0 {
		return "0 second"
	}

	hours := int64(d / time.Hour)
	d %= time.Hour
	minutes := int64(d / time.Minute)
	d %= time.Minute
	seconds := int64(d / time.Second)
	d %= time.Second
	millis := int64(d / time.Millisecond)

	var parts []string
	if hours > 0 {
		parts = append(parts, pluralize(hours, "hour"))
	}
	if minutes > 0 {
		parts = append(parts, pluralize(minutes, "minute"))
	}
	if seconds > 0 {
		parts = append(parts, pluralize(seconds, "second"))
	}
	if millis > 0 {
		parts = append(parts, pluralize(millis, "millisecond"))
	}
	if len(parts) == 0 {
		return "0 second"
	}
	return strings.Join(parts, " ")
}

func pluralize(n int64, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
