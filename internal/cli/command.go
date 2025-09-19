package cli

import viper "github.com/spf13/viper"

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github.com/sijunda/govman/config.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet output (errors only)")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	// Add subcommands
	addCommands()

	// // Remove default `completion` command
	// rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func addCommands() {
	rootCmd.AddCommand(
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
