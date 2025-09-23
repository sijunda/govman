package cli

import (
	"fmt"
	"strings"
	"time"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
	_util "github.com/sijunda/govman/internal/util"
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <version>",
		Short: "📋 Display comprehensive Go version information",
		Long: `Show detailed information about any installed Go version.

📋 Information includes:
  • Version number and release details
  • Complete installation path and directory structure
  • Platform architecture and OS compatibility
  • Installation date, size, and disk usage
  • Binary locations and environment details
  • Release notes and changelog links (when available)

💡 Perfect for debugging installation issues and verifying setups.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			_logger.Verbose("Gathering comprehensive version information for Go %s", version)
			info, err := mgr.Info(version)
			if err != nil {
				helpMsg := fmt.Sprintf("Verify the version is installed with 'govman list', or install it with 'govman install %s'.", version)
				_logger.ErrorWithHelp("Unable to retrieve information for Go %s", helpMsg, version)
				return err
			}

			current, _ := mgr.Current()
			isActive := current == info.Version

			_logger.Info("📋 Go Version Information:")
			_logger.Info(strings.Repeat("═", 60))

			// Basic version info
			activeStatus := "💾 Installed"
			if isActive {
				activeStatus = "✅ Currently Active"
			}
			_logger.Info("📝 Version:          Go %s (%s)", info.Version, activeStatus)

			// Platform information
			_logger.Info("🖥️  Platform:          %s/%s", info.OS, info.Arch)

			// Installation details
			_logger.Info("📁 Installation Path: %s", info.Path)
			_logger.Info("📅 Installed On:      %s", info.InstallDate.Format("Monday, January 2, 2006 at 15:04:05 MST"))
			_logger.Info("💾 Disk Usage:        %s", _util.FormatBytes(info.Size))

			// Additional context
			daysInstalled := int(time.Since(info.InstallDate).Hours() / 24)
			if daysInstalled > 0 {
				_logger.Info("📅 Age:               %d days old", daysInstalled)
			}

			_logger.Info(strings.Repeat("═", 60))

			// Usage suggestions
			if isActive {
				_logger.Info("✅ This version is currently active in your environment")
				_logger.Info("💡 Run 'go version' to verify, or 'go env' to see full environment")
			} else {
				_logger.Info("🔄 Activate this version with: govman use %s", info.Version)
				_logger.Info("📍 Set as default with: govman use %s --default", info.Version)
				_logger.Info("📁 Set for this project: govman use %s --local", info.Version)
			}

			// Show potential issues or recommendations
			if daysInstalled > 180 {
				_logger.Warning("⚠️  This version is over 6 months old - consider updating")
				_logger.Info("🔄 Check for updates with: govman list --remote")
			}

			return nil
		},
	}

	return cmd
}
