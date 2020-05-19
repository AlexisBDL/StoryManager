package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runUser(cmd *cobra.Command, args []string) error {

	fmt.Println(user)

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
