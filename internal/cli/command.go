package cli

import (
	"fmt"
	"os"

	viper "github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github.com/sijunda/govman/config.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet output (errors only)")

	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		// Log the error but continue execution
		fmt.Fprintf(os.Stderr, "Warning: failed to bind verbose flag: %v\n", err)
	}

	if err := viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet")); err != nil {
		// Log the error but continue execution
		fmt.Fprintf(os.Stderr, "Warning: failed to bind quiet flag: %v\n", err)
	}

	// Add subcommands
	addCommands()

	// Remove default `completion` command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func addCommands() {
	rootCmd.AddCommand(
		newInitCmd(),
		newInstallCmd(),
		newUninstallCmd(),
		newUseCmd(),
		newCurrentCmd(),
		newListCmd(),
		newInfoCmd(),
		newCleanCmd(),
		newSelfUpdateCmd(),
	)
}
