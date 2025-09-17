package cli

import (
	cobra "github.com/spf13/cobra"

	_manager "github.com/sijunda/govman/internal/manager"
)

// internal/cli/clean.go
func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean download cache",
		Long:  `Remove cached download files to free up disk space.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())
			return mgr.Clean()
		},
	}

	return cmd
}
