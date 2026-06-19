package cmd

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/Pixie2468/git-branch-explorer/internal/commands"
	"github.com/Pixie2468/git-branch-explorer/internal/tui"
	"github.com/spf13/cobra"
)

var (
	targetPath string
	timeout    time.Duration
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the dir in specific dir with timeout",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Initialize the service using the flag values
		gitSvc := &commands.Service{
			Path:    targetPath,
			Timeout: timeout,
		}

		// Pass the initialized service directly into your model
		p := tea.NewProgram(tui.InitialModel(gitSvc))

		// Run the UI
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running tui: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Bind the flags securely to the variables.
	// StringVarP binds a string. Default is "." (current directory).
	runCmd.Flags().StringVarP(&targetPath, "dir", "d", ".", "Path to the git repository")

	// DurationVarP automatically parses "10s", "1m", etc., into time.Duration.
	runCmd.Flags().DurationVarP(&timeout, "timeout", "t", 0, "Timeout for git commands (e.g., 5s, 1m)")
}
