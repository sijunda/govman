package util

import (
	"fmt"
	"time"
)

// Pre-allocated slice to avoid repeated allocations
var byteSizeUnits = []string{"KB", "MB", "GB", "TB", "PB", "EB"}

// FormatBytes converts a byte count into a human-readable string (KB, MB, GB, ...).
// Parameter size is the number of bytes. Returns a formatted string.
func FormatBytes(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	value := float64(size)
	for i, s := range byteSizeUnits {
		value = value / unit
		if value < unit || i == len(byteSizeUnits)-1 {
			return fmt.Sprintf("%.2f %s", value, s)
		}
	}

	return fmt.Sprintf("%.2f EB+", value)
}

// FormatDuration formats a time.Duration into a concise string (e.g., 45s, 3m12s, 2h05m).
// Parameter d is the duration. Returns a formatted string.
func FormatDuration(d time.Duration) string {
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
