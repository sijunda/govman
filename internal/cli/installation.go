// internal/cli/install.go
package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"

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
				fmt.Printf("\n📦 Installing Go %s...\n", version)
				if err := mgr.Install(version); err != nil {
					errors = append(errors, fmt.Sprintf("Go %s: %v", version, err))
					continue
				}
			}

			if len(errors) > 0 {
				fmt.Printf("\n❌ Some installations failed:\n")
				for _, err := range errors {
					fmt.Printf("  • %s\n", err)
				}
				return fmt.Errorf("failed to install %d version(s)", len(errors))
			}

			fmt.Printf("\n✅ All installations completed successfully!\n")
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

			return mgr.Uninstall(version)
		},
	}

	return cmd
}
