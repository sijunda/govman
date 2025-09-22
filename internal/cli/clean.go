package cli

import (
	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "🧹 Clean download cache and optimize disk usage",
		Long: `Remove cached download files and temporary data to reclaim disk space.

📋 What gets cleaned:
  • Downloaded Go archive files (.tar.gz, .zip)
  • Temporary extraction directories
  • Incomplete or corrupted downloads
  • Obsolete cache metadata and checksums

🔒 Safety guaranteed:
  • Installed Go versions remain completely untouched
  • Your project files and configurations are preserved
  • Only temporary cache files are removed

💡 Run periodically to keep your system clean and optimized.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_logger.Info("🧹 Cleaning download cache and temporary files...")
			_logger.Progress("Scanning cache directories for removable files")

			mgr := _manager.New(getConfig())
			err := mgr.Clean()
			if err != nil {
				_logger.ErrorWithHelp("Unable to clean cache directories", "Verify that ~/.govman/cache exists and you have sufficient permissions to modify it.", "")
				return err
			}

			_logger.Success("✅ Cache cleanup completed successfully")
			_logger.Info("💾 Disk space has been optimized")
			_logger.Info("💡 Your installed Go versions remain untouched and ready to use")
			_logger.Info("🔄 Future downloads will rebuild cache as needed")

			return nil
		},
	}

	return cmd
}
