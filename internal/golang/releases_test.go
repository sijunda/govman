package golang

import (
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := map[string]string{
		"go1.16": "1.16",
		"v1.17":  "1.17",
		"1.18":   "1.18",
	}

	for input, expected := range tests {
		got := normalizeVersion(input)
		if got != expected {
			t.Errorf("normalizeVersion(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestParseVersion(t *testing.T) {
	v := parseVersion("1.17.2-rc3")
	expected := versionParts{
		numbers:    [3]int{1, 17, 2},
		prerelease: "rc3",
	}
	if v != expected {
		t.Errorf("parseVersion failed: got %+v, want %+v", v, expected)
	}

	// Test invalid format
	v2 := parseVersion("invalid")
	if v2.numbers != [3]int{} || v2.prerelease != "" {
		t.Errorf("parseVersion on invalid input should return zero value, got: %+v", v2)
	}
}

func TestComparePrerelease(t *testing.T) {
	tests := []struct {
		pre1     string
		pre2     string
		expected int
	}{
		{"", "", 0},
		{"rc1", "", -1},
		{"", "rc1", 1},
		{"beta2", "rc1", -1},
		{"rc2", "rc1", 1},
		{"alpha3", "alpha2", 1},
		{"beta1", "beta1", 0},
	}

	for _, tt := range tests {
		t.Run(tt.pre1+" vs "+tt.pre2, func(t *testing.T) {
			got := comparePrerelease(tt.pre1, tt.pre2)
			if got != tt.expected {
				t.Errorf("comparePrerelease(%q, %q) = %v, want %v", tt.pre1, tt.pre2, got, tt.expected)
			}
		})
	}
}

func TestGetPrereleaseRank(t *testing.T) {
	tests := map[string]int{
		"rc1":    3,
		"beta2":  2,
		"alpha5": 1,
		"zzz":    0,
	}

	for input, expected := range tests {
		if got := getPrereleaseRank(input); got != expected {
			t.Errorf("getPrereleaseRank(%q) = %v, want %v", input, got, expected)
		}
	}
}

func TestExtractPrereleaseNumber(t *testing.T) {
	tests := map[string]int{
		"rc1":    1,
		"beta2":  2,
		"alpha3": 3,
		"rc":     0,
		"xxx":    0,
	}

	for input, expected := range tests {
		if got := extractPrereleaseNumber(input); got != expected {
			t.Errorf("extractPrereleaseNumber(%q) = %v, want %v", input, got, expected)
		}
	}
}

func TestIsValidVersion(t *testing.T) {
	valid := []string{
		"1.16.0",
		"1.17",
		"1.17.1rc1",
		"1.18beta2",
		"1.19alpha1",
	}
	invalid := []string{
		"",
		"abc",
		"1",
		"1.16..0",
		"v1.16.0",
	}

	for _, v := range valid {
		if !IsValidVersion(v) {
			t.Errorf("IsValidVersion(%q) = false, want true", v)
		}
	}

	for _, v := range invalid {
		if IsValidVersion(v) {
			t.Errorf("IsValidVersion(%q) = true, want false", v)
		}
	}
}

func TestResolveArch(t *testing.T) {
	// Simulate macOS ARM64
	if resolveArch("1.15", "darwin", "arm64") != "amd64" {
		t.Error("Expected darwin/arm64 for Go < 1.16 to resolve to amd64")
	}
	if resolveArch("1.16", "darwin", "arm64") != "arm64" {
		t.Error("Expected darwin/arm64 for Go >= 1.16 to remain arm64")
	}

	// Non-darwin platform
	if resolveArch("1.15", "linux", "amd64") != "amd64" {
		t.Error("Expected linux/amd64 to remain unchanged")
	}
}
