package commands

import (
	"github.com/spf13/cobra"
)

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Manage a story",
}

func init() {
	RootCmd.AddCommand(storyCmd)
}
