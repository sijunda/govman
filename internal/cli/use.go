package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

// getActivationMode returns a human-readable description of the activation mode
func getActivationMode(setDefault, setLocal bool) string {
	if setLocal {
		return "project-local"
	}
	if setDefault {
		return "system-default"
	}
	return "session-only"
}

func newUseCmd() *cobra.Command {
	var (
		setDefault bool
		setLocal   bool
	)

	cmd := &cobra.Command{
		Use:   "use <version>",
		Short: "🔄 Switch between Go versions with flexible activation options",
		Long: `Activate a specific Go version for your development environment.

🎯 Activation Modes:
  • Session-only: Temporary activation for current terminal session
  • System default: Permanent activation across all new sessions
  • Project-local: Version tied to specific project directory

🔍 Smart Features:
  • Automatic verification of version installation
  • Shell integration with PATH management
  • Project-specific .govman-version file support
  • Seamless switching between versions

💡 Examples:
  govman use 1.25.1                 # Session-only activation
  govman use 1.25.1 --default       # Set as system default
  govman use 1.25.1 --local         # Project-specific version`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			// Verify the version is installed first
			if !mgr.IsInstalled(version) {
				_logger.ErrorWithHelp("Go version %s is not installed", "Install it first with 'govman install %s', or check available versions with 'govman list'.", version, version)
				return fmt.Errorf("version %s not installed", version)
			}

			_logger.Verbose("Activating Go %s with mode: %s", version, getActivationMode(setDefault, setLocal))

			// The manager function will print the shell command to stdout for session-only use.
			// We still want our logs to go to stderr.
			err := mgr.Use(version, setDefault, setLocal)
			if err != nil {
				_logger.ErrorWithHelp("Failed to activate Go %s", "Ensure the version is properly installed and you have sufficient permissions.", version)
				return err
			}

			if setLocal {
				_logger.Success("📁 Set Go %s as local version for this project", version)
				_logger.Info("📝 Created/updated .govman-version file in current directory")
				_logger.Info("💡 This version will be used automatically when working in this project")
			} else if setDefault {
				_logger.Success("🏠 Set Go %s as system default version", version)
				_logger.Info("✅ All new terminal sessions will use this version")
				_logger.Info("💡 Current session updated - run 'go version' to verify")
			} else {
				// For session-only use, the manager has printed the necessary command to stdout.
				// We add a confirmation message to stderr.
				_logger.Success("📱 Now using Go %s for this session", version)
				_logger.Info("⚠️  This is temporary - use --default to make it permanent")
				_logger.Info("💡 Run 'go version' to confirm the switch")
			}

			// Show additional context
			info, err := mgr.Info(version)
			if err == nil {
				_logger.Info("💾 Version details: %s/%s, installed %s", info.OS, info.Arch, info.InstallDate.Format("2006-01-02"))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "🏠 Set as system-wide default version (persistent)")
	cmd.Flags().BoolVarP(&setLocal, "local", "l", false, "📁 Set as project-local version (creates .govman-version file)")

	return cmd
}
