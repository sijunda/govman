package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_manager "github.com/sijunda/govman/internal/manager"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [version...]",
		Short: "Install one or more Go versions",
		Long: `Install one or more Go versions from the official Go releases.

Examples:
  govman install latest          # Install latest stable version
  govman install 1.21.5          # Install specific version
  govman install 1.21.5 1.20.12  # Install multiple versions
  govman install 1.22rc1         # Install pre-release version`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			var errors []string
			for _, version := range args {
				if err := mgr.Install(version); err != nil {
					errors = append(errors, fmt.Sprintf("Go %s: %v", version, err))
					continue
				}
			}

			if len(errors) > 0 {
				_logger.ErrorWithHelp("Some installations failed:", "Check the errors above and try installing the versions individually.", "")
				for _, err := range errors {
					_logger.Info("  • %s", err)
				}
				return fmt.Errorf("failed to install %d version(s)", len(errors))
			}

			_logger.Success("All installations completed successfully!")
			return nil
		},
	}

	return cmd
}

func newUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "uninstall <version>",
		Short:   "Uninstall a Go version",
		Long:    `Remove an installed Go version from your system.`,
		Aliases: []string{"remove", "rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			err := mgr.Uninstall(version)
			if err != nil {
				_logger.ErrorWithHelp("Failed to uninstall Go %s", "Make sure the version is installed and not currently active.", version)
				return err
			}
			_logger.Success("Go %s uninstalled successfully", version)
			return nil
		},
	}

	return cmd
}
