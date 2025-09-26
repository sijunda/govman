package cli

import (
	"fmt"
	"os"
	"strings"
	"sync"

	cobra "github.com/spf13/cobra"

	_config "github.com/sijunda/govman/internal/config"
	_version "github.com/sijunda/govman/internal/version"
)

var (
	cfgFile  string
	cfg      *_config.Config
	cfgMutex sync.Mutex
	cfgOnce  sync.Once
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
		"📁 Project-specific versions with .govman-version file support",
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

// Execute runs the root command and shows a banner when no arguments are provided.
func Execute() error {
	// Show banner only when no subcommands or arguments are provided
	// os.Args[0] is always the program name, so we check if there's only that element
	if len(os.Args) <= 1 {
		showBanner()
	}
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
	var initErr error
	cfgOnce.Do(func() {
		var err error
		cfg, err = _config.Load(cfgFile)
		if err != nil {
			initErr = fmt.Errorf("failed to load config: %w", err)
		}
	})
	return initErr
}

func getConfig() *_config.Config {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()
	return cfg
}
