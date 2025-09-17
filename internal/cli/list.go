package cli

import (
	"fmt"
	"path/filepath"

	cobra "github.com/spf13/cobra"

	_manager "github.com/sijunda/govman/internal/manager"
	_util "github.com/sijunda/govman/internal/util"
)

func newListCmd() *cobra.Command {
	var (
		remote     bool
		stableOnly bool
		beta       bool
		pattern    string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List installed Go versions",
		Long:    `List all installed Go versions, with the current version marked.`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := _manager.New(getConfig())

			if remote {
				return listRemoteVersions(mgr, !stableOnly || beta, pattern)
			}

			return listInstalledVersions(mgr)
		},
	}

	cmd.Flags().BoolVarP(&remote, "remote", "r", false, "List available versions for download")
	cmd.Flags().BoolVar(&stableOnly, "stable-only", false, "Show only stable versions (remote only)")
	cmd.Flags().BoolVar(&beta, "beta", false, "Include beta/rc versions (remote only)")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Filter versions by pattern (remote only)")

	return cmd
}

func listInstalledVersions(mgr *_manager.Manager) error {
	versions, err := mgr.ListInstalled()
	if err != nil {
		return fmt.Errorf("failed to list installed versions: %w", err)
	}

	if len(versions) == 0 {
		fmt.Println("No Go versions installed")
		fmt.Println("Run 'govman install latest' to install the latest version")
		return nil
	}

	current, _ := mgr.Current()

	fmt.Println("Installed Go versions:")
	for _, version := range versions {
		marker := "  "
		if version == current {
			marker = "* "
		}

		info, err := mgr.Info(version)
		if err != nil {
			fmt.Printf("%s%s (error getting info)\n", marker, version)
			continue
		}

		size := _util.FormatBytes(info.Size)
		fmt.Printf("%s%s (%s)\n", marker, version, size)
	}

	if current != "" {
		fmt.Printf("\nCurrent: %s\n", current)
	}

	return nil
}

func listRemoteVersions(mgr *_manager.Manager, includeUnstable bool, pattern string) error {
	versions, err := mgr.ListRemote(includeUnstable)
	if err != nil {
		return fmt.Errorf("failed to list remote versions: %w", err)
	}

	// Filter by pattern if provided
	if pattern != "" {
		var filtered []string
		for _, version := range versions {
			if matched, _ := filepath.Match(pattern, version); matched {
				filtered = append(filtered, version)
			}
		}
		versions = filtered
	}

	if len(versions) == 0 {
		fmt.Println("No versions found")
		return nil
	}

	fmt.Println("Available Go versions:")
	for _, version := range versions {
		installed := mgr.IsInstalled(version)
		marker := "  "
		if installed {
			marker = "✓ "
		}
		fmt.Printf("%s%s\n", marker, version)
	}

	return nil
}
