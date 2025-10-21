package cli

import (
	"fmt"
	"os"
	"strings"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

// newRefreshCmd creates the 'refresh' Cobra command to re-evaluate the current directory for a .govman-version file.
// Returns a *cobra.Command whose RunE switches to the local version if present, otherwise to the default; errors if the required version isn't installed.
func newRefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh Go version based on current directory context",
		Long: `Manually trigger version switching based on the current directory.

Purpose:
  • Re-evaluate the current directory for .govman-version files
  • Switch to the appropriate version (local or default)
  • Useful after adding/removing .govman-version files

Examples:
  govman refresh                    # Re-evaluate current directory

Behavior:
  • If .govman-version exists: switch to that version
  • If no .govman-version: switch to default version
  • Equivalent to the auto-switch that happens on 'cd'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			cfg := getConfig()
			filename := cfg.AutoSwitch.ProjectFile
			if data, err := os.ReadFile(filename); err == nil {
				version := strings.TrimSpace(string(data))

				_logger.Info("Found local version file: %s", filename)
				_logger.Info("Switching to Go %s", version)

				if !mgr.IsInstalled(version) {
					helpMsg := fmt.Sprintf("Install it first with 'govman install %s'", version)
					_logger.ErrorWithHelp("Go version %s is not installed", helpMsg, version)
					return fmt.Errorf("version %s not installed", version)
				}

				return mgr.Use(version, false, false)
			}

			_logger.Info("No local version file found")
			_logger.Info("Switching to default Go version")

			return mgr.Use("default", false, false)
		},
	}

	return cmd
}
