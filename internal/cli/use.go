package cli

import (
	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

func newUseCmd() *cobra.Command {
	var (
		setDefault bool
		setLocal   bool
	)

	cmd := &cobra.Command{
		Use:   "use <version>",
		Short: "Switch to a Go version",
		Long: `Switch to a specific Go version.

Examples:
  govman use 1.25.1           # Use version for current session only
  govman use 1.25.1 --default # Set as system default
  govman use 1.25.1 --local   # Set for current project`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			// The manager function will print the shell command to stdout for session-only use.
			// We still want our logs to go to stderr.
			_logger.Info("🐹 Switching to Go %s...", version)

			err := mgr.Use(version, setDefault, setLocal)
			if err != nil {
				_logger.ErrorWithHelp("Failed to switch to Go %s", "Make sure the version is installed. Use 'govman list' to see installed versions.", version)
				return err
			}

			if setLocal {
				_logger.Success("Set local Go version to %s", version)
			} else if setDefault {
				_logger.Success("Set Go %s as default version", version)
			} else {
				// For session-only use, the manager has printed the necessary command to stdout.
				// We add a confirmation message to stderr.
				_logger.Success("Now using Go %s", version)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "Set as system default version")
	cmd.Flags().BoolVarP(&setLocal, "local", "l", false, "Set as project-local version")

	return cmd
}
