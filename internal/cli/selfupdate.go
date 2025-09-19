package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_version "github.com/sijunda/govman/internal/version"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
}

func newSelfUpdateCmd() *cobra.Command {
	var (
		checkOnly  bool
		force      bool
		prerelease bool
	)

	cmd := &cobra.Command{
		Use:   "selfupdate",
		Short: "Update govman to the latest version",
		Long: `Check for and install the latest version of govman.

Examples:
  govman selfupdate              # Update to latest stable version
  govman selfupdate --check      # Check for updates without installing
  govman selfupdate --prerelease # Include prereleases
  govman selfupdate --force      # Force update even if already latest`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSelfUpdate(checkOnly, force, prerelease)
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "check for updates without installing")
	cmd.Flags().BoolVar(&force, "force", false, "force update even if already latest")
	cmd.Flags().BoolVar(&prerelease, "prerelease", false, "include prereleases")

	return cmd
}

func runSelfUpdate(checkOnly, force, prerelease bool) error {
	_logger.Info("Checking for govman updates...")

	_logger.Verbose("Retrieving latest release information")
	latest, err := getLatestRelease(prerelease)
	if err != nil {
		_logger.ErrorWithHelp("Failed to check for updates", "Check your internet connection and try again.", "")
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	current := _version.BuildVersion()
	if current == "dev" {
		_logger.Info("Development version detected, skipping update check")
		return nil
	}

	_logger.Info("Current version: %s", current)
	_logger.Info("Latest version:  %s", latest.TagName)

	if !force && latest.TagName == current {
		_logger.Success("You are already using the latest version!")
		return nil
	}

	if checkOnly {
		if latest.TagName != current {
			_logger.Info("A new version is available: %s", latest.TagName)
			if latest.Body != "" {
				_logger.Info("\nRelease notes:\n%s", latest.Body)
			}
		}
		return nil
	}

	// Find appropriate asset for current platform
	assetName := fmt.Sprintf("govman-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	var downloadURL string
	for _, asset := range latest.Assets {
		if strings.Contains(asset.Name, assetName) {
			downloadURL = asset.DownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	_logger.Download("Downloading %s...", latest.TagName)

	// Download the binary
	_logger.Verbose("Downloading binary")
	resp, err := http.Get(downloadURL)
	if err != nil {
		_logger.ErrorWithHelp("Failed to download binary", "Check your internet connection and try again.", "")
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	// Create a temporary file to save the downloaded binary
	tempFile, err := os.CreateTemp("", "govman-update-*.bin")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the downloaded binary to the temporary file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write binary to temporary file: %w", err)
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Get the path of the current binary
	_logger.Verbose("Getting current binary path")
	currentBinary, err := os.Executable()
	if err != nil {
		_logger.ErrorWithHelp("Failed to get current binary path", "Check if the binary has proper permissions.", "")
		return fmt.Errorf("failed to get current binary path: %w", err)
	}

	// Rename the current binary to a backup
	_logger.Verbose("Creating backup of current binary")
	backupBinary := currentBinary + ".bak"
	if err := os.Rename(currentBinary, backupBinary); err != nil {
		_logger.ErrorWithHelp("Failed to create backup of current binary", "Check if you have permission to modify the binary directory.", "")
		return fmt.Errorf("failed to rename current binary to backup: %w", err)
	}

	// Move the downloaded binary to the current binary path
	_logger.Verbose("Installing new binary")
	if err := os.Rename(tempFile.Name(), currentBinary); err != nil {
		// Restore the backup if the move fails
		_logger.Warning("Failed to install new binary, restoring backup")
		if err := os.Rename(backupBinary, currentBinary); err != nil {
			_logger.ErrorWithHelp("Failed to restore backup binary", "You may need to manually restore the binary from the backup file.", "")
			return fmt.Errorf("failed to restore backup binary: %w", err)
		}
		return fmt.Errorf("failed to move downloaded binary to current binary path: %w", err)
	}

	// Set the executable permission for the new binary
	_logger.Verbose("Setting executable permissions")
	if err := os.Chmod(currentBinary, 0755); err != nil {
		_logger.ErrorWithHelp("Failed to set executable permissions", "You may need to manually set executable permissions on the binary.", "")
		return fmt.Errorf("failed to set executable permission for new binary: %w", err)
	}

	_logger.Success("Update completed successfully!")
	return nil
}

func getLatestRelease(includePrerelease bool) (*GitHubRelease, error) {
	url := "https://api.github.com/repos/sijunda/govman/releases/latest"
	if includePrerelease {
		url = "https://api.github.com/repos/sijunda/govman/releases?per_page=1"
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if includePrerelease {
		var releases []GitHubRelease
		if err := json.Unmarshal(body, &releases); err != nil {
			return nil, err
		}
		if len(releases) == 0 {
			return nil, fmt.Errorf("no releases found")
		}
		// Find the first prerelease or stable release
		for _, release := range releases {
			if includePrerelease || !release.Prerelease {
				return &release, nil
			}
		}
		return &releases[0], nil
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}
