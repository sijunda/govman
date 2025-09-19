package cli

import (
	"fmt"
	"strings"

	cobra "github.com/spf13/cobra"

	_config "github.com/sijunda/govman/internal/config"
	_version "github.com/sijunda/govman/internal/version"
)

var (
	cfgFile string
	cfg     *_config.Config
)

var rootCmd = &cobra.Command{
	Use:     "govman",
	Short:   "🚀 Go Version Manager - Install and manage multiple Go versions",
	Long:    createLongDescription(),
	Version: _version.BuildVersion(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

func createLongDescription() string {
	features := []string{
		"⚡ Lightning-fast installation and switching between Go versions",
		"🎯 Zero configuration - works out of the box, no setup required",
		"📁 Project-specific versions with .go-version file support",
		"🚫 No admin/sudo required - fully userspace installation",
		"💾 Intelligent caching with offline mode support",
		"📦 Parallel downloads with automatic resume on failure",
		"🌍 Cross-platform support (Windows, macOS, Linux, ARM)",
		"🧹 Built-in cleanup tools to manage disk space efficiently",
	}

	var sb strings.Builder
	sb.WriteString("\n🎯 Key Features:\n")
	for _, feature := range features {
		sb.WriteString(fmt.Sprintf("  %s\n", feature))
	}
	return sb.String()
}

func Execute() error {
	showBanner()
	return rootCmd.Execute()
}

func showBanner() {
	fmt.Println()
	banner := `
	 ██████╗  ██████╗ ██╗   ██╗███╗   ███╗ █████╗ ███╗   ██╗
	██╔════╝ ██╔═══██╗██║   ██║████╗ ████║██╔══██╗████╗  ██║
	██║  ███╗██║   ██║██║   ██║██╔████╔██║███████║██╔██╗ ██║
	██║   ██║██║   ██║╚██╗ ██╔╝██║╚██╔╝██║██╔══██║██║╚██╗██║
	╚██████╔╝╚██████╔╝ ╚████╔╝ ██║ ╚═╝ ██║██║  ██║██║ ╚████║
	 ╚═════╝  ╚═════╝   ╚═══╝  ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝`

	lines := strings.Split(banner, "\n")

	const (
		color = "\033[38;5;75m"
		bold  = "\033[1m"
		reset = "\033[0m"
	)

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Println(color + bold + line + reset)
		}
	}
	fmt.Println()
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
