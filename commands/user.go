package commands

import (
	"fmt"

	"github.com/AlexisBDL/StoryManager/config"
	"github.com/spf13/cobra"
)

func runUser(cmd *cobra.Command, args []string) error {
	cfg := config.NewResolver()

	fmt.Println(cfg.GetUserString())

	return nil
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show user of database",
	Args:  cobra.ExactArgs(0),
	RunE:  runUser,
}

func init() {
	RootCmd.AddCommand(userCmd)
}
