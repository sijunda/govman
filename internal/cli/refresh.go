package cli

import (
	"fmt"
	"os"
	"strings"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

func newRefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "🔄 Refresh Go version based on current directory context",
		Long: `Manually trigger version switching based on the current directory.

🎯 Purpose:
  • Re-evaluate the current directory for .govman-version files
  • Switch to the appropriate version (local or default)
  • Useful after adding/removing .govman-version files

💡 Examples:
  govman refresh                    # Re-evaluate current directory

🔍 Behavior:
  • If .govman-version exists: switch to that version
  • If no .govman-version: switch to default version
  • Equivalent to the auto-switch that happens on 'cd'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			// Check for local version file
			cfg := getConfig()
			filename := cfg.AutoSwitch.ProjectFile
			if data, err := os.ReadFile(filename); err == nil {
				// Local version file exists
				version := strings.TrimSpace(string(data)) // Remove whitespace/newlines

				_logger.Info("📁 Found local version file: %s", filename)
				_logger.Info("🔄 Switching to Go %s", version)

				// Verify the version is installed
				if !mgr.IsInstalled(version) {
					helpMsg := fmt.Sprintf("Install it first with 'govman install %s'", version)
					_logger.ErrorWithHelp("Go version %s is not installed", helpMsg, version)
					return fmt.Errorf("version %s not installed", version)
				}

				// Use session-only mode (same as auto-switch behavior)
				return mgr.Use(version, false, false)
			} else {
				// No local version file, switch to default
				_logger.Info("📂 No local version file found")
				_logger.Info("🔄 Switching to default Go version")

				return mgr.Use("default", false, false)
			}
		},
	}

	return cmd
}