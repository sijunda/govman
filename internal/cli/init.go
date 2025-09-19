package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"

	_logger "github.com/sijunda/govman/internal/logger"
	_shell "github.com/sijunda/govman/internal/shell"
)

func newInitCmd() *cobra.Command {
	var (
		force     bool
		shellName string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize shell integration",
		Long: `Set up shell integration for automatic Go version switching.
This command adds govman to your shell configuration and enables
auto-switching based on .govman-version files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var sh _shell.Shell

			if shellName != "" {
				sh = getShellByName(shellName)
				if sh == nil {
					return fmt.Errorf("unsupported shell: %s", shellName)
				}
			} else {
				sh = _shell.Detect()
			}

			cfg := getConfig()
			binPath := cfg.GetBinPath()

			_logger.Info("🔧 Initializing %s integration...", sh.Name())

			_logger.Step("Setting up shell integration")
			if err := _shell.InitializeShell(sh, binPath, force); err != nil {
				_logger.ErrorWithHelp("Failed to initialize shell integration", "Check if you have permission to modify your shell configuration file.", "")
				return err
			}

			_logger.Success("Shell integration initialized!")
			_logger.Info("📝 Please restart your shell or run: source %s", sh.ConfigFile())

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-initialization")
	cmd.Flags().StringVar(&shellName, "shell", "", "Target shell (bash, zsh, fish, powershell)")

	return cmd
}

func getShellByName(name string) _shell.Shell {
	switch name {
	case "bash":
		return &_shell.BashShell{}
	case "zsh":
		return &_shell.ZshShell{}
	case "fish":
		return &_shell.FishShell{}
	case "powershell", "pwsh":
		return &_shell.PowerShell{}
	default:
		return nil
	}
}
