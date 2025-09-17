package version

import (
	"fmt"
	"runtime"
)

var (
	// These will be set by ldflags during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
	buildBy = "unknown"
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
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuildBy:   buildBy,
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
