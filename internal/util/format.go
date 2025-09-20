package util

import (
	"fmt"
	"time"
)

// FormatBytes formats a byte size into a human-readable string with appropriate units (B, KB, MB, GB, etc.)
func FormatBytes(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	sizes := []string{"KB", "MB", "GB", "TB", "PB", "EB"}

	value := float64(size)
	for i, s := range sizes {
		value = value / unit
		if value < unit || i == len(sizes)-1 {
			return fmt.Sprintf("%.2f %s", value, s)
		}
	}

	// Fallback (should never be hit)
	return fmt.Sprintf("%.2f EB+", value)
}

// FormatDuration formats a duration into a human-readable string.
// It shows seconds for durations less than a minute, minutes and seconds
// for durations less than an hour, and hours and minutes for durations
// of an hour or more.
func FormatDuration(d time.Duration) string {
	// Handle negative durations by formatting the absolute value
	// and adding a minus sign prefix
	if d < 0 {
		return "-" + FormatDuration(-d)
	}

	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}
