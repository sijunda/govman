package util

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	testCases := []struct {
		name     string
		size     int64
		expected string
	}{
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
