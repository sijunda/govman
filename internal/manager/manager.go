package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	_config "github.com/sijunda/govman/internal/config"
	_downloader "github.com/sijunda/govman/internal/downloader"
	_golang "github.com/sijunda/govman/internal/golang"
	_logger "github.com/sijunda/govman/internal/logger"
	_shell "github.com/sijunda/govman/internal/shell"
	_symlink "github.com/sijunda/govman/internal/symlink"
)

type Manager struct {
	config     *_config.Config
	downloader *_downloader.Downloader
	shell      _shell.Shell
}

func New(cfg *_config.Config) *Manager {
	return &Manager{
		config:     cfg,
		downloader: _downloader.New(cfg),
		shell:      _shell.Detect(),
	}
}

// Install installs a Go version
func (m *Manager) Install(version string) error {
	// Resolve version (latest, etc.)
	timer := _logger.StartTimer("version resolution")
	resolvedVersion, err := m.resolveVersion(version)
	if err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to resolve version %s: %w", version, err)
	}
	_logger.StopTimer(timer)

	// Check if already installed
	_logger.InternalProgress("Checking if version is already installed")
	if m.IsInstalled(resolvedVersion) {
		return fmt.Errorf("go version %s is already installed", resolvedVersion)
	}

	// Download and install
	_logger.Info("Installing Go %s...", resolvedVersion)

	timer = _logger.StartTimer("download URL retrieval")
	downloadURL, err := _golang.GetDownloadURLWithConfig(resolvedVersion,
		m.config.GoReleases.APIURL,
		m.config.GoReleases.CacheExpiry,
		m.config.GoReleases.DownloadURL)
	if err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to get download URL: %w", err)
	}
	_logger.StopTimer(timer)

	installDir := m.config.GetVersionDir(resolvedVersion)
	timer = _logger.StartTimer("download and installation")
	if err := m.downloader.Download(downloadURL, installDir, resolvedVersion); err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to download and install: %w", err)
	}
	_logger.StopTimer(timer)

	_logger.Success("Go %s installed successfully", resolvedVersion)
	return nil
}

// Uninstall removes a Go version
func (m *Manager) Uninstall(version string) error {
	_logger.InternalProgress("Checking if version is installed")
	if !m.IsInstalled(version) {
		return fmt.Errorf("go version %s is not installed", version)
	}

	// Check if it's the current version
	_logger.InternalProgress("Checking if version is currently active")
	current, err := m.Current()
	if err == nil && current == version {
		return fmt.Errorf("cannot uninstall currently active version %s", version)
	}

	installDir := m.config.GetVersionDir(version)
	_logger.InternalProgress("Removing installation directory: %s", installDir)
	timer := _logger.StartTimer("uninstallation")
	if err := os.RemoveAll(installDir); err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to remove installation directory: %w", err)
	}
	_logger.StopTimer(timer)

	_logger.Success("Go %s uninstalled successfully", version)
	return nil
}

// Use switches to a Go version
func (m *Manager) Use(version string, setDefault, setLocal bool) error {
	// Handle special "default" version
	if version == "default" {
		// Get the default version from the symlink
		defaultVersion, err := m.CurrentGlobal()
		if err != nil {
			return fmt.Errorf("failed to get default version: %w", err)
		}
		version = defaultVersion
	} else {
		_logger.InternalProgress("Checking if version is installed")
		if !m.IsInstalled(version) {
			return fmt.Errorf("go version %s is not installed. Run 'govman install %s' first", version, version)
		}
	}

	// Set local version (project-specific)
	if setLocal {
		_logger.InternalProgress("Setting local version")
		if err := m.setLocalVersion(version); err != nil {
			return fmt.Errorf("failed to set local version: %w", err)
		}
		return nil
	}

	// Set as default (persistent)
	if setDefault {
		_logger.InternalProgress("Creating symlink")
		timer := _logger.StartTimer("symlink creation")
		if err := m.createSymlink(version); err != nil {
			_logger.StopTimer(timer)
			return fmt.Errorf("failed to create symlink: %w", err)
		}
		_logger.StopTimer(timer)

		_logger.InternalProgress("Setting as default version")
		m.config.DefaultVersion = version
		timer = _logger.StartTimer("saving configuration")
		if err := m.config.Save(); err != nil {
			_logger.StopTimer(timer)
			return fmt.Errorf("failed to save default version: %w", err)
		}
		_logger.StopTimer(timer)
		return nil
	}

	// Session-only use: directly modify PATH in-process if possible (e.g., in interactive shell)
	// But in general, we'll call a shell command to export it
	versionBinPath := filepath.Join(m.config.GetVersionDir(version), "bin")

	// Inject to os.Environ for current process
	os.Setenv("PATH", versionBinPath+string(os.PathListSeparator)+os.Getenv("PATH"))

	// Print the shell command to stdout so the shell wrapper can evaluate it
	// Detect shell type and output appropriate command
	var shellCmd string
	shell := m.shell.Name()
	if shell == "fish" {
		shellCmd = fmt.Sprintf("set -gx PATH \"%s\" $PATH", versionBinPath)
	} else if shell == "powershell" {
		shellCmd = fmt.Sprintf("$env:PATH = \"%s;\" + $env:PATH", versionBinPath)
	} else {
		// Default to bash/zsh format
		shellCmd = fmt.Sprintf("export PATH=\"%s:$PATH\"", versionBinPath)
	}
	fmt.Println(shellCmd)

	// Optionally, export to shell for confirmation (could also be omitted)
	_logger.Verbose("PATH updated to include: %s", versionBinPath)

	return nil
}

// Current returns the currently active Go version
// It prioritizes local version files over global settings
func (m *Manager) Current() (string, error) {
	// Check for local version first (project-specific .govman-version file)
	if localVersion := m.getLocalVersion(); localVersion != "" {
		// Validate that the local version is actually installed
		if !m.IsInstalled(localVersion) {
			return "", fmt.Errorf("local version %s specified in %s is not installed - run 'govman install %s' to install it",
				localVersion, m.config.AutoSwitch.ProjectFile, localVersion)
		}
		return localVersion, nil
	}

	// Check if there's a session-only version active by checking the 'go' command in PATH
	// This is important because 'govman use' with session-only mode modifies the PATH
	// but doesn't change the symlink, so the actual active version might be different
	sessionVersion, err := m.getCurrentSessionVersion()
	if err == nil && sessionVersion != "" {
		return sessionVersion, nil
	}

	version, err := m.CurrentGlobal()
	if err != nil {
		return "", err
	}

	return version, nil
}

// CurrentGlobal returns the globally active Go version from the symlink
func (m *Manager) CurrentGlobal() (string, error) {
	// Check for symlink to determine global active version
	symlinkPath := m.config.GetCurrentSymlink()

	// Check if symlink exists
	linkInfo, err := os.Lstat(symlinkPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No symlink exists - check if we have a default version configured
			if m.config.DefaultVersion != "" {
				if m.IsInstalled(m.config.DefaultVersion) {
					return "", fmt.Errorf("no active Go version found - default version %s is configured but symlink is missing. Run 'govman use %s' to activate it",
						m.config.DefaultVersion, m.config.DefaultVersion)
				} else {
					return "", fmt.Errorf("no active Go version found - default version %s is configured but not installed. Run 'govman install %s' first, then 'govman use %s'",
						m.config.DefaultVersion, m.config.DefaultVersion, m.config.DefaultVersion)
				}
			}
			return "", fmt.Errorf("no Go version is currently active - no symlink found at %s and no default version configured. Install a version with 'govman install <version>' and activate it with 'govman use <version>'",
				symlinkPath)
		}
		return "", fmt.Errorf("failed to check symlink at %s: %w - this may indicate a permissions issue or corrupted installation",
			symlinkPath, err)
	}

	// Verify it's actually a symlink
	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("expected symlink at %s but found %s instead - this may indicate a corrupted govman installation. Try running 'govman use <version>' to recreate the symlink",
			symlinkPath, linkInfo.Mode().Type().String())
	}

	// Read the symlink target
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink target from %s: %w - the symlink may be corrupted",
			symlinkPath, err)
	}

	// Extract version from the target path
	// Expected path format: /path/to/versions/go1.21.0/bin/go
	// We need to extract "1.25.1" from this path
	targetDir := filepath.Dir(target)      // Remove /go from the end
	targetDir = filepath.Dir(targetDir)    // Remove /bin from the end
	versionDir := filepath.Base(targetDir) // Get go1.21.0

	if !strings.HasPrefix(versionDir, "go") {
		return "", fmt.Errorf("invalid symlink target format: expected version directory to start with 'go' but found %s - the symlink may be corrupted. Target path: %s",
			versionDir, target)
	}

	version := versionDir[2:] // Remove "go" prefix to get "1.25.1"

	// Validate the extracted version
	if version == "" {
		return "", fmt.Errorf("could not extract version from symlink target %s - the symlink may be corrupted", target)
	}

	// Verify the version directory still exists
	expectedVersionDir := m.config.GetVersionDir(version)
	if _, err := os.Stat(expectedVersionDir); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("symlink points to Go %s but installation directory %s no longer exists - the installation may have been manually deleted. Run 'govman install %s' to reinstall",
				version, expectedVersionDir, version)
		}
		return "", fmt.Errorf("failed to verify installation directory %s for Go %s: %w",
			expectedVersionDir, version, err)
	}

	// Verify the actual Go executable exists and is functional
	goExecutable := filepath.Join(expectedVersionDir, "bin", "go")
	// On Windows, the executable has a .exe extension
	if runtime.GOOS == "windows" {
		goExecutable += ".exe"
	}
	if _, err := os.Stat(goExecutable); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("go %s installation appears corrupted - executable not found at %s. Try reinstalling with 'govman install %s'",
				version, goExecutable, version)
		}
		return "", fmt.Errorf("failed to verify Go executable at %s for version %s: %w",
			goExecutable, version, err)
	}

	return version, nil
}

// ListInstalled returns all installed Go versions
func (m *Manager) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(m.config.InstallDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read install directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "go") {
			version := entry.Name()[2:] // Remove "go" prefix
			versions = append(versions, version)
		}
	}

	// Sort versions
	sort.Slice(versions, func(i, j int) bool {
		return _golang.CompareVersions(versions[i], versions[j]) > 0
	})

	return versions, nil
}

// ListRemote returns all available Go versions for download
func (m *Manager) ListRemote(includeUnstable bool) ([]string, error) {
	return _golang.GetAvailableVersionsWithConfig(includeUnstable,
		m.config.GoReleases.APIURL,
		m.config.GoReleases.CacheExpiry)
}

// IsInstalled checks if a version is installed
func (m *Manager) IsInstalled(version string) bool {
	installDir := m.config.GetVersionDir(version)
	_, err := os.Stat(installDir)
	return err == nil
}

// Info returns information about a version
func (m *Manager) Info(version string) (*_golang.VersionInfo, error) {
	if !m.IsInstalled(version) {
		return nil, fmt.Errorf("go version %s is not installed", version)
	}

	installDir := m.config.GetVersionDir(version)
	return _golang.GetVersionInfo(installDir)
}

// Clean removes cached files
func (m *Manager) Clean() error {
	if err := os.RemoveAll(m.config.CacheDir); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	// Recreate cache directory
	if err := os.MkdirAll(m.config.CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate cache directory: %w", err)
	}

	_logger.Success("Cache cleaned successfully")
	return nil
}

// resolveVersion resolves aliases like "latest" or partial versions like "1.24"
func (m *Manager) resolveVersion(version string) (string, error) {
	if version == "latest" {
		versions, err := m.ListRemote(false) // false for stable only
		if err != nil {
			return "", err
		}
		if len(versions) == 0 {
			return "", fmt.Errorf("no stable versions available")
		}
		return versions[0], nil
	}

	// Check if it's a partial version like "1.24"
	if strings.Count(version, ".") == 1 {
		// Get all versions including unstable to have the complete set
		versions, err := m.ListRemote(true) // true to include unstable versions
		if err != nil {
			return "", err
		}

		prefix := version + "."
		for _, v := range versions {
			if strings.HasPrefix(v, prefix) {
				// The list is sorted newest first, so the first match is the latest patch version.
				return v, nil
			}
		}
	}

	return version, nil
}

// createSymlink creates or updates the symlink to the specified version
func (m *Manager) createSymlink(version string) error {
	// targetDir is the path to the version's root directory (e.g., /Users/sijunda/.govman/versions/go1.25.1)
	versionRoot := m.config.GetVersionDir(version)
	// The actual Go executable is inside the 'bin' directory within the version's root
	goExecutablePath := filepath.Join(versionRoot, "bin", "go")

	// On Windows, the executable has a .exe extension
	if runtime.GOOS == "windows" {
		goExecutablePath += ".exe"
	}

	symlinkPath := m.config.GetCurrentSymlink() // This gets the path to the symlink (e.g., /Users/sijunda/.govman/bin/go)

	// On Windows, the symlink should also have a .exe extension
	if runtime.GOOS == "windows" {
		symlinkPath += ".exe"
	}

	// Create bin directory if it doesn't exist
	binDir := m.config.GetBinPath()
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Remove existing symlink
	os.Remove(symlinkPath)

	// Create new symlink
	if err := _symlink.Create(goExecutablePath, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// setLocalVersion creates a .govman-version file in the current directory
func (m *Manager) setLocalVersion(version string) error {
	filename := m.config.AutoSwitch.ProjectFile
	return os.WriteFile(filename, []byte(version), 0644)
}

// getLocalVersion reads version from .govman-version file if it exists
func (m *Manager) getLocalVersion() string {
	filename := m.config.AutoSwitch.ProjectFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// DefaultVersion returns the default Go version from the config
func (m *Manager) DefaultVersion() string {
	return m.config.DefaultVersion
}

// GetDefaultVersionFromSymlink returns the version pointed to by the symlink
func (m *Manager) GetDefaultVersionFromSymlink() (string, error) {
	return m.CurrentGlobal()
}

// getCurrentSessionVersion checks if there's a session-only version active by running 'go version'
// and parsing the output to determine which version is currently in the PATH
func (m *Manager) getCurrentSessionVersion() (string, error) {
	// Execute 'go version' to see which version is currently active in the PATH
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'go version': %w", err)
	}

	// Parse the output to extract the version
	// Expected format: "go version go1.25.1 darwin/amd64" or similar
	versionStr := strings.TrimSpace(string(output))
	parts := strings.Split(versionStr, " ")
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected 'go version' output format: %s", versionStr)
	}

	version := strings.TrimPrefix(parts[2], "go")
	if version == "" {
		return "", fmt.Errorf("could not extract version from 'go version' output: %s", versionStr)
	}

	return version, nil
}
