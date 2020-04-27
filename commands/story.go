package commands

import (
	"github.com/spf13/cobra"
)

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Show, Create, delete or set a story",
}

func init() {
	RootCmd.AddCommand(storyCmd)
}
