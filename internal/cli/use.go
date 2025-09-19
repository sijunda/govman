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
		shellEval  bool
	)

	cmd := &cobra.Command{
		Use:   "use <version>",
		Short: "Switch to a Go version",
		Long: `Switch to a specific Go version.

Examples:
  govman use 1.21.5           # Use version for current session
  govman use 1.21.5 --default # Set as system default
  govman use 1.21.5 --local   # Set for current project`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			cfg := getConfig()
			mgr := _manager.New(cfg)

			// For shell evaluation, we want to output only the shell command
			if shellEval {
				// For shell evaluation, suppress all logging
				// and handle errors by writing to stderr directly
			} else if setDefault || setLocal {
				_logger.Info("🐹 Switching to Go %s...", version)
			}

			// Use the new method with logging control
			logOutput := !shellEval
			err := mgr.UseWithOptions(version, setDefault, setLocal, logOutput)
			if err != nil {
				if shellEval {
					// For shell evaluation, output error to stderr
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					return err
				} else {
					_logger.ErrorWithHelp("Failed to switch to Go %s", "Make sure the version is installed. Use 'govman list' to see installed versions.", version)
					return err
				}
			}

			if setLocal {
				_logger.Success("Set local Go version to %s", version)
			} else if setDefault {
				_logger.Success("Set Go %s as default version", version)
			} else {
				// For session-only changes or shell evaluation, output shell commands to modify PATH
				versionDir := cfg.GetVersionDir(version)
				binPath := filepath.Join(versionDir, "bin")
				fmt.Printf("export PATH=\"%s:$PATH\"\n", binPath)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "Set as system default version")
	cmd.Flags().BoolVarP(&setLocal, "local", "l", false, "Set as project-local version")
	cmd.Flags().BoolVar(&shellEval, "shell-eval", false, "Output shell command for evaluation")
	// Hide the shell-eval flag from help as it's for internal use
	cmd.Flags().MarkHidden("shell-eval")

	return cmd
}
