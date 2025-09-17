package util

import "fmt"

// FormatBytes formats a byte size into a human-readable string with appropriate units (B, KB, MB, GB, etc.)
func FormatBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	kb := float64(size) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%.2f KB", kb)
	}

	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%.2f MB", mb)
	}

	gb := mb / 1024
	return fmt.Sprintf("%.2f GB", gb)
}
