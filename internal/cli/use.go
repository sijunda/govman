package cli

import (
	"fmt"
	"os"
	"path/filepath"

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

This command has three modes of operation:

1. Session-only (default): Activates the version for the current terminal session only.
		 Requires evaluating the output: eval "$(govman use 1.21.5)"
		 When the terminal is closed, the version change is lost.

2. Persistent (--default): Sets the version as the system default.
		 The version remains active even after closing and reopening the terminal.

3. Project-local (--local): Sets the version for the current project only.
		 Creates a .govman-version file that automatically activates the version
		 when you're in that directory.

Examples:
		eval "$(govman use 1.21.5)"  # Use version for current session
		govman use 1.21.5 --default  # Set as system default (persistent)
		govman use 1.21.5 --local    # Set for current project (local)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			cfg := getConfig()
			mgr := _manager.New(cfg)

			// Handle different modes
			if setDefault || setLocal {
				_logger.Info("🐹 Switching to Go %s...", version)

				// Use the new method with logging enabled
				err := mgr.UseWithOptions(version, setDefault, setLocal, true)
				if err != nil {
					_logger.ErrorWithHelp("Failed to switch to Go %s", "Make sure the version is installed. Use 'govman list' to see installed versions.", version)
					return err
				}

				if setLocal {
					_logger.Success("Set local Go version to %s", version)
				} else if setDefault {
					_logger.Success("Set Go %s as default version", version)
				}
			} else {
				// Session-only usage - check if version is installed first
				if !mgr.IsInstalled(version) {
					_logger.ErrorWithHelp("Go version %s is not installed", "Install the version first with 'govman install %s'", version, version)
					return fmt.Errorf("go version %s is not installed", version)
				}

				// For session-only changes, output shell commands to modify PATH
				versionDir := cfg.GetVersionDir(version)
				binPath := filepath.Join(versionDir, "bin")
				fmt.Printf("export PATH=\"%s:$PATH\"\n", binPath)
				fmt.Fprintf(os.Stderr, "📝 To activate this version in your current session, run:\n")
				fmt.Fprintf(os.Stderr, "   eval \"$(govman use %s)\"\n", version)
				fmt.Fprintf(os.Stderr, "   or manually run: export PATH=\"%s:$PATH\"\n", binPath)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "Set as system default version (persistent)")
	cmd.Flags().BoolVarP(&setLocal, "local", "l", false, "Set as project-local version")

	return cmd
}
