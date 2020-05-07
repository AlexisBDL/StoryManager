package commands

import (
	"github.com/spf13/cobra"
)

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "MAnage a story",
}

func init() {
	RootCmd.AddCommand(storyCmd)
}
