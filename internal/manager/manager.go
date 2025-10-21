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

// New constructs a Manager with the provided configuration.
// It initializes a downloader and detects the user's shell.
func New(cfg *_config.Config) *Manager {
	return &Manager{
		config:     cfg,
		downloader: _downloader.New(cfg),
		shell:      _shell.Detect(),
	}
}

// Install downloads and installs the specified Go version.
// version may be an exact string or "latest". Returns an error if resolution, download, or installation fails.
func (m *Manager) Install(version string) error {
	timer := _logger.StartTimer("version resolution")
	resolvedVersion, err := m.resolveVersion(version)
	if err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to resolve version %s: %w", version, err)
	}
	_logger.StopTimer(timer)

	_logger.InternalProgress("Checking if version is already installed")
	if m.IsInstalled(resolvedVersion) {
		return fmt.Errorf("go version %s is already installed", resolvedVersion)
	}

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

// Uninstall removes an installed Go version.
// Returns an error if the version is not installed, is active, or removal fails.
func (m *Manager) Uninstall(version string) error {
	_logger.InternalProgress("Checking if version is installed")
	if !m.IsInstalled(version) {
		return fmt.Errorf("go version %s is not installed", version)
	}

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

// Use activates a Go version for the current session, as default, or for the local project.
// setDefault sets it globally; setLocal writes a project version file. Returns an error if activation fails.
func (m *Manager) Use(version string, setDefault, setLocal bool) error {
	if version == "default" {
		defaultVersion, err := m.CurrentGlobal()
		if err != nil {
			return fmt.Errorf("failed to get default version: %w", err)
		}
		version = defaultVersion
	} else {
		// Validate version is installed
		_logger.InternalProgress("Checking if version is installed")
		if !m.IsInstalled(version) {
			return fmt.Errorf("go version %s is not installed. Run 'govman install %s' first", version, version)
		}
	}

	// Apply the version based on scope
	switch {
	case setLocal:
		_logger.InternalProgress("Setting local version for project")
		if err := m.setLocalVersion(version); err != nil {
			return fmt.Errorf("failed to set local version: %w", err)
		}
		_logger.Success("Set Go %s as local version for this project", version)

	case setDefault:
		_logger.InternalProgress("Setting as system default version")

		// Update config
		m.config.DefaultVersion = version
		if err := m.config.Save(); err != nil {
			_logger.Warning("Failed to save default version to config: %v", err)
		}

		// Create symlink
		_logger.InternalProgress("Creating symlink for Go %s", version)
		timer := _logger.StartTimer("symlink creation")
		if err := m.createSymlink(version); err != nil {
			_logger.StopTimer(timer)
			return fmt.Errorf("failed to create symlink: %w", err)
		}
		_logger.StopTimer(timer)

	default:
		// Session-only, no additional action needed
	}

	// Update PATH
	versionBinPath := filepath.Join(m.config.GetVersionDir(version), "bin")
	return m.shell.ExecutePathCommand(versionBinPath)
}

// Current returns the currently active Go version, checking session, local project, or global symlink.
// Returns the version string or an error if none is active or validation fails.
func (m *Manager) Current() (string, error) {
	sessionVersion, err := m.getCurrentSessionVersion()
	if err == nil && sessionVersion != "" {
		if !m.IsInstalled(sessionVersion) {
			_logger.Warning("Session version %s is active but not managed by GOVMAN", sessionVersion)
		}

		return sessionVersion, nil
	}

	if localVersion := m.getLocalVersion(); localVersion != "" {
		if !m.IsInstalled(localVersion) {
			return "", fmt.Errorf("local version %s specified in %s is not installed - run 'govman install %s' to install it",
				localVersion, m.config.AutoSwitch.ProjectFile, localVersion)
		}

		return localVersion, nil
	}

	version, err := m.CurrentGlobal()
	if err != nil {
		return "", err
	}

	return version, nil
}

// CurrentGlobal resolves the active global version from the symlink and validates installation integrity.
// Returns the version or an error for missing/corrupt symlink or installation.
func (m *Manager) CurrentGlobal() (string, error) {
	symlinkPath := m.config.GetCurrentSymlink()

	linkInfo, err := os.Lstat(symlinkPath)
	if err != nil {
		if os.IsNotExist(err) {
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

	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("expected symlink at %s but found %s instead - this may indicate a corrupted govman installation. Try running 'govman use <version>' to recreate the symlink",
			symlinkPath, linkInfo.Mode().Type().String())
	}

	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink target from %s: %w - the symlink may be corrupted",
			symlinkPath, err)
	}

	targetDir := filepath.Dir(target)
	targetDir = filepath.Dir(targetDir)
	versionDir := filepath.Base(targetDir)

	if !strings.HasPrefix(versionDir, "go") {
		return "", fmt.Errorf("invalid symlink target format: expected version directory to start with 'go' but found %s - the symlink may be corrupted. Target path: %s",
			versionDir, target)
	}

	version := versionDir[2:]
	if version == "" {
		return "", fmt.Errorf("could not extract version from symlink target %s - the symlink may be corrupted", target)
	}

	expectedVersionDir := m.config.GetVersionDir(version)
	if _, err := os.Stat(expectedVersionDir); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("symlink points to Go %s but installation directory %s no longer exists - the installation may have been manually deleted. Run 'govman install %s' to reinstall",
				version, expectedVersionDir, version)
		}

		return "", fmt.Errorf("failed to verify installation directory %s for Go %s: %w",
			expectedVersionDir, version, err)
	}

	goExecutable := filepath.Join(expectedVersionDir, "bin", "go")

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

// ListInstalled returns installed Go versions sorted in descending order.
// Returns the slice of versions or an error if the install directory cannot be read.
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
			version := entry.Name()[2:]
			versions = append(versions, version)
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return _golang.CompareVersions(versions[i], versions[j]) > 0
	})

	return versions, nil
}

// ListRemote fetches available remote Go versions.
// includeUnstable controls inclusion of beta/rc versions. Returns the list or an error.
func (m *Manager) ListRemote(includeUnstable bool) ([]string, error) {
	return _golang.GetAvailableVersionsWithConfig(includeUnstable,
		m.config.GoReleases.APIURL,
		m.config.GoReleases.CacheExpiry)
}

// IsInstalled reports whether a given version is installed by checking its directory.
// Returns true if installed; false otherwise.
func (m *Manager) IsInstalled(version string) bool {
	installDir := m.config.GetVersionDir(version)
	_, err := os.Stat(installDir)

	return err == nil
}

// Info returns metadata about an installed version.
// Returns VersionInfo or an error if the version is not installed or info retrieval fails.
func (m *Manager) Info(version string) (*_golang.VersionInfo, error) {
	if !m.IsInstalled(version) {
		return nil, fmt.Errorf("go version %s is not installed", version)
	}

	installDir := m.config.GetVersionDir(version)
	return _golang.GetVersionInfo(installDir)
}

// Clean removes and recreates the cache directory.
// Returns an error if cleanup fails; nil on success.
func (m *Manager) Clean() error {
	if err := os.RemoveAll(m.config.CacheDir); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	if err := os.MkdirAll(m.config.CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate cache directory: %w", err)
	}

	_logger.Success("Cache cleaned successfully")
	return nil
}

// resolveVersion resolves aliases and partial versions to a concrete version.
// "latest" becomes the newest stable; "major.minor" expands to the latest patch. Returns the resolved version or an error.
func (m *Manager) resolveVersion(version string) (string, error) {
	if version == "latest" {
		versions, err := m.ListRemote(false)
		if err != nil {
			return "", err
		}

		if len(versions) == 0 {
			return "", fmt.Errorf("no stable versions available")
		}

		return versions[0], nil
	}

	if strings.Count(version, ".") == 1 {
		versions, err := m.ListRemote(true)
		if err != nil {
			return "", err
		}

		prefix := version + "."
		for _, v := range versions {
			if strings.HasPrefix(v, prefix) {
				return v, nil
			}
		}
	}

	return version, nil
}

// createSymlink creates/replaces the global "go" symlink targeting the selected version's binary.
// Returns an error if directory creation or symlink operation fails.
func (m *Manager) createSymlink(version string) error {
	versionRoot := m.config.GetVersionDir(version)

	goExecutablePath := filepath.Join(versionRoot, "bin", "go")

	if runtime.GOOS == "windows" {
		goExecutablePath += ".exe"
	}

	symlinkPath := m.config.GetCurrentSymlink()

	if runtime.GOOS == "windows" {
		symlinkPath += ".exe"
	}

	binDir := m.config.GetBinPath()
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Remove the old symlink if it exists
	if err := os.Remove(symlinkPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing symlink: %w", err)
	}

	if err := _symlink.Create(goExecutablePath, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// setLocalVersion writes the project's autoswitch file with the specified version.
// Returns an error if the file write fails.
func (m *Manager) setLocalVersion(version string) error {
	filename := m.config.AutoSwitch.ProjectFile
	return os.WriteFile(filename, []byte(version), 0644)
}

// getLocalVersion reads the project's autoswitch file and returns the local version.
// Returns an empty string if the file does not exist or cannot be read.
func (m *Manager) getLocalVersion() string {
	filename := m.config.AutoSwitch.ProjectFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// DefaultVersion returns the configured default version string.
func (m *Manager) DefaultVersion() string {
	return m.config.DefaultVersion
}

// GetDefaultVersionFromSymlink returns the active/default version by reading the global symlink.
// It delegates to CurrentGlobal and returns its result.
func (m *Manager) GetDefaultVersionFromSymlink() (string, error) {
	return m.CurrentGlobal()
}

// CurrentActivationMethod returns the activation method for the currently active Go version.
// Returns "session-only", "project-local", or "system-default" based on how the current version is activated.
func (m *Manager) CurrentActivationMethod() string {
	sessionVersion, err := m.getCurrentSessionVersion()
	if err == nil && sessionVersion != "" {
		if localVersion := m.getLocalVersion(); localVersion != "" && localVersion == sessionVersion {
			return "project-local"
		}

		globalVersion, err := m.CurrentGlobal()
		if err == nil && globalVersion == sessionVersion {
			return "system-default"
		}

		return "session-only"
	}

	if localVersion := m.getLocalVersion(); localVersion != "" {
		return "project-local"
	}

	return "system-default"
}

// getCurrentSessionVersion executes "go version" and parses the active version.
// Returns the version string or an error if command execution or parsing fails.
func (m *Manager) getCurrentSessionVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'go version': %w", err)
	}

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
