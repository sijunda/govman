package cli

import (
	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
	_util "github.com/sijunda/govman/internal/util"
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <version>",
		Short: "Show information about a Go version",
		Long:  `Display detailed information about an installed Go version.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			_logger.Verbose("Retrieving version information")
			info, err := mgr.Info(version)
			if err != nil {
				_logger.ErrorWithHelp("Failed to get information for Go version %s", "Make sure the version is installed. Use 'govman list' to see installed versions.", version)
				return err
			}

			_logger.Info("Go Version: %s", info.Version)
			_logger.Info("Install Path: %s", info.Path)
			_logger.Info("Platform: %s/%s", info.OS, info.Arch)
			_logger.Info("Install Date: %s", info.InstallDate.Format("2006-01-02 15:04:05"))
			_logger.Info("Size: %s", _util.FormatBytes(info.Size))

			return nil
		},
	}

	return cmd
}
