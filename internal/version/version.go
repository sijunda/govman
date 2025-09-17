package version

import (
	"fmt"
	"runtime"
)

var (
	// These will be set by ldflags during build
	// IMPORTANT: Must be exported (uppercase) to be set by linker
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuildBy = "unknown"
)

// Info represents version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	BuildBy   string `json:"buildBy"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		BuildBy:   BuildBy,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// BuildVersion returns a version string for cobra
func BuildVersion() string {
	info := Get()
	if info.Version == "dev" {
		return fmt.Sprintf("%s-%s", info.Version, info.Commit)
	}
	return info.Version
}
