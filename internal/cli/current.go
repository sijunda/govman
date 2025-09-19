package cli

import (
	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

func newCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show current Go version",
		Long:  `Display the currently active Go version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			_logger.Step("Retrieving current version")
			current, err := mgr.Current()
			if err != nil {
				_logger.Info("No Go version is currently active")
				return nil
			}

			_logger.Info("Current Go version: %s", current)
			return nil
		},
	}

	return cmd
}
