package cli

import (
	cobra "github.com/spf13/cobra"

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
  govman use 1.21.5           # Use version for current session
  govman use 1.21.5 --default # Set as system default
  govman use 1.21.5 --local   # Set for current project`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			return mgr.Use(version, setDefault, setLocal)
		},
	}

	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "Set as system default version")
	cmd.Flags().BoolVarP(&setLocal, "local", "l", false, "Set as project-local version")

	return cmd
}
