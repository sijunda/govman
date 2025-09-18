// internal/cli/selfupdate.go
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
	fmt.Println("Checking for govman updates...")

	latest, err := getLatestRelease(prerelease)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	current := _version.BuildVersion()
	if current == "dev" {
		fmt.Println("Development version detected, skipping update check")
		return nil
	}

	fmt.Printf("Current version: %s\n", current)
	fmt.Printf("Latest version:  %s\n", latest.TagName)

	if !force && latest.TagName == current {
		fmt.Println("You are already using the latest version!")
		return nil
	}

	if checkOnly {
		if latest.TagName != current {
			fmt.Printf("A new version is available: %s\n", latest.TagName)
			if latest.Body != "" {
				fmt.Printf("\nRelease notes:\n%s\n", latest.Body)
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

	fmt.Printf("Downloading %s...\n", latest.TagName)

	// Download the binary
	resp, err := http.Get(downloadURL)
	if err != nil {
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
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current binary path: %w", err)
	}

	// Rename the current binary to a backup
	backupBinary := currentBinary + ".bak"
	if err := os.Rename(currentBinary, backupBinary); err != nil {
		return fmt.Errorf("failed to rename current binary to backup: %w", err)
	}

	// Move the downloaded binary to the current binary path
	if err := os.Rename(tempFile.Name(), currentBinary); err != nil {
		// Restore the backup if the move fails
		if err := os.Rename(backupBinary, currentBinary); err != nil {
			return fmt.Errorf("failed to restore backup binary: %w", err)
		}
		return fmt.Errorf("failed to move downloaded binary to current binary path: %w", err)
	}

	// Set the executable permission for the new binary
	if err := os.Chmod(currentBinary, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission for new binary: %w", err)
	}

	fmt.Println("Update completed successfully!")
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
