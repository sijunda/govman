package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"

	_manager "github.com/sijunda/govman/internal/manager"
	_util "github.com/sijunda/govman/internal/util"
)

func newCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show current Go version",
		Long:  `Display the currently active Go version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			current, err := mgr.Current()
			if err != nil {
				return fmt.Errorf("no Go version is currently active")
			}

			fmt.Println(current)
			return nil
		},
	}

	return cmd
}

// internal/cli/info.go
func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <version>",
		Short: "Show information about a Go version",
		Long:  `Display detailed information about an installed Go version.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			info, err := mgr.Info(version)
			if err != nil {
				return err
			}

			fmt.Printf("Go Version: %s\n", info.Version)
			fmt.Printf("Install Path: %s\n", info.Path)
			fmt.Printf("Platform: %s/%s\n", info.OS, info.Arch)
			fmt.Printf("Install Date: %s\n", info.InstallDate.Format("2006-01-02 15:04:05"))
			fmt.Printf("Size: %s\n", _util.FormatBytes(info.Size))

			return nil
		},
	}

	return cmd
}
