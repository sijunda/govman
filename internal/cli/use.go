package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
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
  govman use 1.21.5           # Use version for current session only
  govman use 1.21.5 --default # Set as system default
  govman use 1.21.5 --local   # Set for current project`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			mgr := _manager.New(getConfig())

			// Check if both flags are set (invalid)
			if setDefault && setLocal {
				_logger.ErrorWithHelp("Invalid flags", "Cannot use --default and --local together")
				return fmt.Errorf("cannot use --default and --local flags together")
			}

			_logger.Info("🐹 Switching to Go %s...", version)
			
			// For session-only usage (no flags), we handle it differently
			if !setDefault && !setLocal {
				return handleSessionOnlyUse(mgr, version)
			}

			// For persistent usage, use the existing logic
			err := mgr.Use(version, setDefault, setLocal)
			if err != nil {
				_logger.ErrorWithHelp("Failed to switch to Go %s", "Make sure the version is installed. Use 'govman list' to see installed versions.", version)
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "Set as system default version")
	cmd.Flags().BoolVarP(&setLocal, "local", "l", false, "Set as project-local version")

	return cmd
}

func handleSessionOnlyUse(mgr *_manager.Manager, version string) error {
	// Check if version is installed
	if !mgr.IsInstalled(version) {
		_logger.ErrorWithHelp("Go version %s is not installed", "Run 'govman install %s' first", version, version)
		return fmt.Errorf("go version %s is not installed", version)
	}

	config := getConfig()
	versionRoot := config.GetVersionDir(version)
	goBinPath := filepath.Join(versionRoot, "bin")

	// Get current PATH
	currentPath := os.Getenv("PATH")
	
	// Remove any existing govman paths from PATH to avoid conflicts
	cleanPath := removeGovmanPaths(currentPath, config.GetBinPath())
	
	// Prepend the version-specific bin path
	newPath := goBinPath + string(os.PathListSeparator) + cleanPath

	// Print shell commands for the user to execute
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash" // fallback
	}

	shellName := filepath.Base(shell)
	
	_logger.Success("To use Go %s in this session, run:", version)
	
	switch shellName {
	case "zsh", "bash":
		fmt.Printf("export PATH=\"%s\"\n", newPath)
	case "fish":
		fmt.Printf("set -x PATH \"%s\"\n", newPath)
	case "csh", "tcsh":
		fmt.Printf("setenv PATH \"%s\"\n", newPath)
	default:
		fmt.Printf("export PATH=\"%s\"\n", newPath)
	}
	
	_logger.Info("This will only affect your current terminal session.")
	_logger.Info("Close and reopen your terminal to revert to the previous Go version.")
	
	return nil
}

// removeGovmanPaths removes govman bin paths from the PATH string
func removeGovmanPaths(path, govmanBinPath string) string {
	if path == "" {
		return ""
	}
	
	paths := filepath.SplitList(path)
	var cleanPaths []string
	
	for _, p := range paths {
		// Skip paths that are the govman bin directory
		if p != govmanBinPath {
			cleanPaths = append(cleanPaths, p)
		}
	}
	
	return filepath.Join(cleanPaths...)
}