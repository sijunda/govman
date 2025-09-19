package cli

import (
	"fmt"
	"path/filepath"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
	_util "github.com/sijunda/govman/internal/util"
)

func newListCmd() *cobra.Command {
	var (
		remote     bool
		stableOnly bool
		beta       bool
		pattern    string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List installed Go versions",
		Long:    `List all installed Go versions, with the current version marked.`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			if remote {
				return listRemoteVersions(mgr, !stableOnly || beta, pattern)
			}

			return listInstalledVersions(mgr)
		},
	}

	cmd.Flags().BoolVarP(&remote, "remote", "r", false, "List available versions for download")
	cmd.Flags().BoolVar(&stableOnly, "stable-only", false, "Show only stable versions (remote only)")
	cmd.Flags().BoolVar(&beta, "beta", false, "Include beta/rc versions (remote only)")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Filter versions by pattern (remote only)")

	return cmd
}

func listInstalledVersions(mgr *_manager.Manager) error {
	_logger.Verbose("Retrieving installed versions")
	versions, err := mgr.ListInstalled()
	if err != nil {
		_logger.ErrorWithHelp("Failed to list installed versions", "Check if the installation directory is accessible.", "")
		return fmt.Errorf("failed to list installed versions: %w", err)
	}

	if len(versions) == 0 {
		_logger.Info("No Go versions installed")
		_logger.Info("Run 'govman install latest' to install the latest version")
		return nil
	}

	current, _ := mgr.Current()

	_logger.Info("Installed Go versions:")
	for _, version := range versions {
		marker := "  "
		if version == current {
			marker = "* "
		}

		info, err := mgr.Info(version)
		if err != nil {
			_logger.Info("%s%s (error getting info)", marker, version)
			continue
		}

		size := _util.FormatBytes(info.Size)
		_logger.Info("%s%s (%s)", marker, version, size)
	}

	if current != "" {
		_logger.Info("\nCurrent: %s", current)
	}

	return nil
}

func listRemoteVersions(mgr *_manager.Manager, includeUnstable bool, pattern string) error {
	_logger.Verbose("Retrieving remote versions")
	versions, err := mgr.ListRemote(includeUnstable)
	if err != nil {
		_logger.ErrorWithHelp("Failed to list remote versions", "Check your internet connection and try again.", "")
		return fmt.Errorf("failed to list remote versions: %w", err)
	}

	// Filter by pattern if provided
	if pattern != "" {
		var filtered []string
		for _, version := range versions {
			if matched, _ := filepath.Match(pattern, version); matched {
				filtered = append(filtered, version)
			}
		}
		versions = filtered
	}

	if len(versions) == 0 {
		_logger.Info("No versions found")
		return nil
	}

	_logger.Info("Available Go versions:")
	for _, version := range versions {
		installed := mgr.IsInstalled(version)
		marker := "  "
		if installed {
			marker = "✓ "
		}
		_logger.Info("%s%s", marker, version)
	}

	return nil
}
