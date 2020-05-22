// Package commands contains the CLI commands
package commands

import (
	"os"

	"github.com/AlexisBDL/StoryManager/config"

	"github.com/spf13/cobra"
)

const rootCommandName = "StoryManager"

// Init
var (
	cfg  *config.Resolver
	user string
)

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
	cfg = config.NewResolver() //config default db "Stories"
	user = cfg.GetUserString()
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
