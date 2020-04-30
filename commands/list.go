package commands

import (
	"fmt"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func runListStory(cmd *cobra.Command, args []string) error {
	cfg := config.NewResolver() //config default db "Stories"
	db, err := cfg.GetDatabase("")
	d.CheckError(err)
	defer db.Close()

	db.Datasets().IterAll(func(k, v types.Value) {
		fmt.Println(k)
	})

	return nil
}

var listStoryCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all stories in database Stories.",
	Args:  cobra.ExactArgs(0),
	RunE:  runListStory,
}

func init() {
	RootCmd.AddCommand(listStoryCmd)
}
