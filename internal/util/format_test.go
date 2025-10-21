package util

import (
	"testing"
	"time"
)

func TestFormatBytes(t *testing.T) {
	testCases := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "Zero bytes",
			size:     0,
			expected: "0 B",
		},
		{
			name:     "Bytes less than 1 KB",
			size:     512,
			expected: "512 B",
		},
		{
			name:     "Exactly 1 KB",
			size:     1024,
			expected: "1.00 KB",
		},
		{
			name:     "Multiple KB",
			size:     2048,
			expected: "2.00 KB",
		},
		{
			name:     "Just below 1 MB",
			size:     1024*1024 - 1,
			expected: "1024.00 KB",
		},
		{
			name:     "Exactly 1 MB",
			size:     1024 * 1024,
			expected: "1.00 MB",
		},
		{
			name:     "Multiple MB",
			size:     5 * 1024 * 1024,
			expected: "5.00 MB",
		},
		{
			name:     "Just below 1 GB",
			size:     1024*1024*1024 - 1,
			expected: "1024.00 MB",
		},
		{
			name:     "Exactly 1 GB",
			size:     1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "Multiple GB",
			size:     3 * 1024 * 1024 * 1024,
			expected: "3.00 GB",
		},
		{
			name:     "TB size",
			size:     2 * 1024 * 1024 * 1024 * 1024,
			expected: "2.00 TB",
		},
		{
			name:     "PB size",
			size:     3 * 1024 * 1024 * 1024 * 1024 * 1024,
			expected: "3.00 PB",
		},
		{
			name:     "EB size (largest unit)",
			size:     4 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024,
			expected: "4.00 EB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatBytes(tc.size)
			if result != tc.expected {
				t.Errorf("FormatBytes(%d) = %q; want %q", tc.size, result, tc.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "Zero duration",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "Seconds only",
			duration: 30 * time.Second,
			expected: "30s",
		},
		{
			name:     "Just under 1 minute",
			duration: 59 * time.Second,
			expected: "59s",
		},
		{
			name:     "Exactly 1 minute",
			duration: time.Minute,
			expected: "1m0s",
		},
		{
			name:     "Minutes and seconds",
			duration: 2*time.Minute + 30*time.Second,
			expected: "2m30s",
		},
		{
			name:     "Just under 1 hour",
			duration: 59*time.Minute + 59*time.Second,
			expected: "59m59s",
		},
		{
			name:     "Exactly 1 hour",
			duration: time.Hour,
			expected: "1h0m",
		},
		{
			name:     "Hours and minutes",
			duration: 2*time.Hour + 30*time.Minute,
			expected: "2h30m",
		},
		{
			name:     "Hours, minutes, and seconds (seconds ignored)",
			duration: 3*time.Hour + 45*time.Minute + 30*time.Second,
			expected: "3h45m",
		},
		{
			name:     "Large duration",
			duration: 25*time.Hour + 15*time.Minute,
			expected: "25h15m",
		},
		{
			name:     "Negative duration - seconds",
			duration: -30 * time.Second,
			expected: "-30s",
		},
		{
			name:     "Negative duration - minutes",
			duration: -2*time.Minute - 30*time.Second,
			expected: "-2m30s",
		},
		{
			name:     "Negative duration - hours",
			duration: -3*time.Hour - 45*time.Minute,
			expected: "-3h45m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatDuration(tc.duration)
			if result != tc.expected {
				t.Errorf("FormatDuration(%v) = %q; want %q", tc.duration, result, tc.expected)
			}
		})
	}
}
