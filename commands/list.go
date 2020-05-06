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

	var ls []string
	db.Datasets().IterAll(func(k, v types.Value) {
		ls = append(ls, fmt.Sprintln(k))
	})
	//modifier list pour trouver open et close
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
