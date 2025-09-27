package golang

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	releasesCache []Release
	cacheMutex    sync.RWMutex
	cacheExpiry   time.Time
)

const (
	GoDownloadURLTemplate = "%s"
)

var (
	// Default values that can be overridden
	defaultGoReleasesAPI = "https://go.dev/dl/?mode=json&include=all"
	defaultCacheDuration = 10 * time.Minute
	defaultGoDownloadURL = "https://go.dev/dl/%s"
)

type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Sha256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
}

type VersionInfo struct {
	Version     string
	Path        string
	OS          string
	Arch        string
	InstallDate time.Time
	Size        int64
}

// GetAvailableVersions fetches all available Go versions
func GetAvailableVersions(includeUnstable bool) ([]string, error) {
	return GetAvailableVersionsWithConfig(includeUnstable, defaultGoReleasesAPI, defaultCacheDuration)
}

func GetAvailableVersionsWithConfig(includeUnstable bool, apiURL string, cacheDuration time.Duration) ([]string, error) {
	releases, err := fetchReleasesWithConfig(apiURL, cacheDuration)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, release := range releases {
		if !includeUnstable && !release.Stable {
			continue
		}
		// Remove "go" prefix for consistency
		version := strings.TrimPrefix(release.Version, "go")
		versions = append(versions, version)
	}

	// Sort versions (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return CompareVersions(versions[i], versions[j]) > 0
	})

	return versions, nil
}

// GetDownloadURL returns the download URL for a specific version
func GetDownloadURL(version string) (string, error) {
	return GetDownloadURLWithConfig(version, defaultGoReleasesAPI, defaultCacheDuration, defaultGoDownloadURL)
}

func GetDownloadURLWithConfig(version string, apiURL string, cacheDuration time.Duration, downloadURL string) (string, error) {
	releases, err := fetchReleasesWithConfig(apiURL, cacheDuration)
	if err != nil {
		return "", err
	}

	targetVersion := "go" + version
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Handle older Go versions that might not have a direct darwin/arm64 build
	// but have a compatible one.
	resolvedArch := resolveArch(version, goos, goarch)

	for _, release := range releases {
		if release.Version != targetVersion {
			continue
		}

		for _, file := range release.Files {
			if file.OS == goos && file.Arch == resolvedArch && file.Kind == "archive" {
				return fmt.Sprintf(downloadURL, file.Filename), nil
			}
		}
	}

	return "", fmt.Errorf("no download available for Go %s on %s/%s", version, goos, goarch)
}

// resolveArch returns the correct architecture for a given Go version.
// For example, Go versions older than 1.16 don't have a darwin/arm64 build,
// but the darwin/amd64 build can be used on Apple Silicon with Rosetta 2.
func resolveArch(version, goos, goarch string) string {
	// As of Go 1.16, official darwin/arm64 builds are available.
	// For older versions, we can fall back to the amd64 build on Apple Silicon.
	if goos == "darwin" && goarch == "arm64" {
		if CompareVersions(version, "1.16") < 0 {
			return "amd64"
		}
	}
	return goarch
}

// GetFileInfo returns file information for a specific version
func GetFileInfo(version string) (*File, error) {
	return GetFileInfoWithConfig(version, defaultGoReleasesAPI, defaultCacheDuration)
}

func GetFileInfoWithConfig(version string, apiURL string, cacheDuration time.Duration) (*File, error) {
	releases, err := fetchReleasesWithConfig(apiURL, cacheDuration)
	if err != nil {
		return nil, err
	}

	targetVersion := "go" + version
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	resolvedArch := resolveArch(version, goos, goarch)

	for _, release := range releases {
		if release.Version != targetVersion {
			continue
		}

		for _, file := range release.Files {
			if file.OS == goos && file.Arch == resolvedArch && file.Kind == "archive" {
				return &file, nil
			}
		}
	}

	return nil, fmt.Errorf("no file info available for Go %s on %s/%s", version, goos, goarch)
}

// GetVersionInfo returns information about an installed version
func GetVersionInfo(installPath string) (*VersionInfo, error) {
	// Check if go binary exists
	goBinary := filepath.Join(installPath, "bin", "go")
	if runtime.GOOS == "windows" {
		goBinary += ".exe"
	}

	stat, err := os.Stat(goBinary)
	if err != nil {
		return nil, fmt.Errorf("go binary not found in %s", installPath)
	}

	// Extract version from path
	version := filepath.Base(installPath)
	version = strings.TrimPrefix(version, "go")

	// Calculate directory size
	size, err := getDirSize(installPath)
	if err != nil {
		size = 0 // Continue without size info
	}

	return &VersionInfo{
		Version:     version,
		Path:        installPath,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstallDate: stat.ModTime(),
		Size:        size,
	}, nil
}

// CompareVersions compares two version strings
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func CompareVersions(v1, v2 string) int {
	// Quick equality check
	if v1 == v2 {
		return 0
	}

	// Normalize versions (remove prefixes, handle rc/beta)
	v1 = normalizeVersion(v1)
	v2 = normalizeVersion(v2)

	// Check again after normalization
	if v1 == v2 {
		return 0
	}

	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	// Compare major, minor, patch
	for i := 0; i < 3; i++ {
		if parts1.numbers[i] > parts2.numbers[i] {
			return 1
		} else if parts1.numbers[i] < parts2.numbers[i] {
			return -1
		}
	}

	// Compare pre-release (stable > rc > beta > alpha)
	return comparePrerelease(parts1.prerelease, parts2.prerelease)
}

// IsValidVersion checks if a version string is valid
func IsValidVersion(version string) bool {
	// Basic version pattern: x.y.z with optional pre-release
	pattern := `^\d+\.\d+(?:\.\d+)?(?:-?(?:rc|beta|alpha)\d*)?$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

type versionParts struct {
	numbers    [3]int // major, minor, patch
	prerelease string
}

func parseVersion(version string) versionParts {
	var parts versionParts

	// Extract pre-release suffix with improved regex
	re := regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?(?:-?(rc\d+|beta\d+|alpha\d+))?$`)
	matches := re.FindStringSubmatch(version)

	if len(matches) == 0 {
		return parts // Invalid version
	}

	// Parse major version
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			parts.numbers[0] = num
		}
	}

	// Parse minor version
	if len(matches) > 2 {
		if num, err := strconv.Atoi(matches[2]); err == nil {
			parts.numbers[1] = num
		}
	}

	// Parse patch version (optional)
	if len(matches) > 3 && matches[3] != "" {
		if num, err := strconv.Atoi(matches[3]); err == nil {
			parts.numbers[2] = num
		}
	}

	// Extract pre-release (if present)
	if len(matches) > 4 && matches[4] != "" {
		parts.prerelease = matches[4]
	}

	return parts
}

func normalizeVersion(version string) string {
	// Remove common prefixes
	version = strings.TrimPrefix(version, "go")
	version = strings.TrimPrefix(version, "v")
	return version
}

func comparePrerelease(pre1, pre2 string) int {
	if pre1 == "" && pre2 == "" {
		return 0
	}
	if pre1 == "" {
		return 1 // Stable > pre-release
	}
	if pre2 == "" {
		return -1
	}

	// Both are pre-releases, compare by type and number
	rank1 := getPrereleaseRank(pre1)
	rank2 := getPrereleaseRank(pre2)

	if rank1 != rank2 {
		return rank1 - rank2
	}

	// Same type, compare numbers
	num1 := extractPrereleaseNumber(pre1)
	num2 := extractPrereleaseNumber(pre2)

	if num1 > num2 {
		return 1
	} else if num1 < num2 {
		return -1
	}

	return 0
}

func getPrereleaseRank(prerelease string) int {
	if strings.HasPrefix(prerelease, "rc") {
		return 3
	} else if strings.HasPrefix(prerelease, "beta") {
		return 2
	} else if strings.HasPrefix(prerelease, "alpha") {
		return 1
	}
	return 0
}

func extractPrereleaseNumber(prerelease string) int {
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(prerelease)
	if num, err := strconv.Atoi(match); err == nil {
		return num
	}
	return 0
}

func fetchReleases() ([]Release, error) {
	return fetchReleasesWithConfig(defaultGoReleasesAPI, defaultCacheDuration)
}

func fetchReleasesWithConfig(apiURL string, cacheDuration time.Duration) ([]Release, error) {
	// Check cache first
	cacheMutex.RLock()
	if time.Now().Before(cacheExpiry) && releasesCache != nil {
		defer cacheMutex.RUnlock()
		return releasesCache, nil
	}
	cacheMutex.RUnlock()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Fetch from API
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, fmt.Errorf("failed to parse releases: %w", err)
	}

	// Update cache
	cacheMutex.Lock()
	releasesCache = releases
	cacheExpiry = time.Now().Add(cacheDuration)
	cacheMutex.Unlock()

	return releases, nil
}

func getDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			// Continue walking even if some files can't be accessed
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// ClearReleasesCache clears the releases cache
func ClearReleasesCache() {
	cacheMutex.Lock()
	releasesCache = nil
	cacheExpiry = time.Time{}
	cacheMutex.Unlock()
}
