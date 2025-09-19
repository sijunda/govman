package cli

import (
	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean download cache",
		Long:  `Remove cached download files to free up disk space.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_logger.Info("🧹 Cleaning download cache...")
			mgr := _manager.New(getConfig())
			err := mgr.Clean()
			if err != nil {
				_logger.ErrorWithHelp("Failed to clean cache", "Check if you have permission to access the cache directory.", "")
				return err
			}
			_logger.Success("Cache cleaned successfully")
			return nil
		},
	}

	return cmd
}
