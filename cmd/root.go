package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gbx",
	Short: "A Git Branch Explorer TUI",
}

func Execute() error {
	return rootCmd.Execute()
}
