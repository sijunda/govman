package cli

import (
	"fmt"
	"strings"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
	_util "github.com/sijunda/govman/internal/util"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [version...]",
		Short: "🚀 Install Go versions with intelligent download management",
		Long: `Download and install one or more Go versions from official releases.

🎯 Features:
  • Lightning-fast parallel downloads with resume capability
  • Automatic integrity verification and checksum validation
  • Smart caching to avoid re-downloading existing archives
  • Support for latest, stable, and pre-release versions
  • Batch installation with detailed progress tracking
  • Automatic cleanup of temporary files on completion

💡 Examples:
  govman install latest              # Latest stable release
  govman install 1.25.1              # Specific version
  govman install 1.25.1 1.20.12      # Multiple versions
  govman install 1.22rc1             # Pre-release version`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			_logger.Info("🚀 Starting installation of %d Go version(s)...", len(args))
			_logger.Progress("Preparing downloads and verifying version availability")

			var errors []string
			var successful []string
			for i, version := range args {
				_logger.Info("💾 [%d/%d] Installing Go %s...", i+1, len(args), version)
				if err := mgr.Install(version); err != nil {
					errors = append(errors, fmt.Sprintf("❌ Go %s: %v", version, err))
					_logger.Warning("Failed to install Go %s: %v", version, err)
					continue
				}
				successful = append(successful, version)
				_logger.Success("Successfully installed Go %s", version)
			}

			_logger.Info(strings.Repeat("─", 50))

			if len(successful) > 0 {
				_logger.Success("✅ Successfully installed %d version(s):", len(successful))
				for _, version := range successful {
					_logger.Info("  • Go %s", version)
				}
			}

			if len(errors) > 0 {
				_logger.ErrorWithHelp("Failed to install %d version(s):", "Review the errors below and try installing problematic versions individually for more details.", len(errors))
				for _, err := range errors {
					_logger.Info("  %s", err)
				}
				_logger.Info("💡 Common solutions:")
				_logger.Info("  • Check your internet connection")
				_logger.Info("  • Verify version exists with 'govman list --remote'")
				_logger.Info("  • Try again with verbose mode: govman install <version> --verbose")
				return fmt.Errorf("failed to install %d version(s)", len(errors))
			}

			if len(successful) > 0 {
				_logger.Success("🎉 All installations completed successfully!")
				if len(successful) == 1 {
					_logger.Info("💡 Activate it with: govman use %s", successful[0])
				} else {
					_logger.Info("💡 List all versions: govman list")
					_logger.Info("💡 Activate any version: govman use <version>")
				}
			}

			return nil
		},
	}

	return cmd
}

func newUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall <version>",
		Short: "🗑️ Safely remove Go versions with cleanup",
		Long: `Completely remove an installed Go version from your system.

🔒 Safety features:
  • Prevents removal of currently active versions
  • Confirms version exists before attempting removal
  • Complete cleanup of binaries and associated files
  • Automatic recalculation of disk space
  • Preserves other installed versions safely

💡 The uninstalled version will no longer appear in 'govman list'.`,
		Aliases: []string{"remove", "rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			// Check if version is currently active
			current, _ := mgr.Current()
			if current == version {
				_logger.ErrorWithHelp("Cannot uninstall currently active Go version %s", "Switch to a different version first with 'govman use <other-version>', then try uninstalling again.", version)
				return fmt.Errorf("cannot uninstall active version")
			}

			// Get info before uninstalling for better reporting
			info, err := mgr.Info(version)
			if err != nil {
				_logger.ErrorWithHelp("Go version %s is not installed or information is unavailable", "Use 'govman list' to see all installed versions.", version)
				return err
			}

			_logger.Info("🗑️ Uninstalling Go %s...", version)
			_logger.Progress("Removing installation directory and associated files")

			err = mgr.Uninstall(version)
			if err != nil {
				_logger.ErrorWithHelp("Failed to uninstall Go %s", "Ensure no processes are using this Go installation and you have sufficient permissions.", version)
				return err
			}

			_logger.Success("✅ Successfully uninstalled Go %s", version)
			_logger.Info("💾 Freed up %s of disk space", _util.FormatBytes(info.Size))
			_logger.Info("💡 View remaining versions with: govman list")

			return nil
		},
	}

	return cmd
}
