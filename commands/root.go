// Package commands contains the CLI commands
package commands

import (
	"os"

	"github.com/spf13/cobra"
)

const rootCommandName = "StoryManager"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   rootCommandName,
	Short: "A User Story manager",

	// For the root command, force the execution of the PreRun
	// even if we just display the help. This is to make sure that we check
	// the repository and give the user early feedback.
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
	},

	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
