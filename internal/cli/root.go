package cli

import (
	"fmt"

	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"

	_config "github.com/sijunda/govman/internal/config"
	_version "github.com/sijunda/govman/internal/version"
)

var (
	cfgFile string
	cfg     *_config.Config
)

var rootCmd = &cobra.Command{
	Use:   "govman",
	Short: "Go Version Manager - Install and manage multiple Go versions",
	Long: `GOVMAN is a cross-platform Go version manager that allows you to 
install, manage, and switch between multiple Go versions effortlessly.

Features:
• Install and switch between multiple Go versions
• Project-specific version support via .govman-version
• Cross-platform support (Windows, macOS, Linux)
• Automatic shell integration
• Fast parallel downloads with resume capability
• Package manager integration`,
	Version: _version.BuildVersion(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github.com/sijunda/govman/config.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet output (errors only)")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	// Add subcommands
	addCommands()
}

func addCommands() {
	rootCmd.AddCommand(
		newInstallCmd(),
		newListCmd(),
		newUseCmd(),
		newUninstallCmd(),
		newCurrentCmd(),
		newInfoCmd(),
		newInitCmd(),
		newCleanCmd(),
		newSelfUpdateCmd(),
	)
}

func initConfig() error {
	if cfg != nil {
		return nil // Already initialized
	}

	var err error
	cfg, err = _config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}

func getConfig() *_config.Config {
	return cfg
}
