package commands

import (
	"fmt"

	"github.com/AlexisBDL/StoryManager/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/outputpager"

	"github.com/spf13/cobra"
)

func runShowStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	cfg := config.NewResolver() //config default db "Stories"
	db, ds, err := cfg.GetPath(ID)
	d.CheckError(err)
	defer db.Close()

	if ds == nil {
		fmt.Printf("Story %s not found in database\n", ID)
		return nil
	}

	pgr := outputpager.Start()
	defer pgr.Stop()

	types.WriteEncodedValue(pgr.Writer, ds)
	fmt.Fprintln(pgr.Writer)

	return nil
}

var showStoryCmd = &cobra.Command{
	Use:   "show <ID>",
	Short: "show a story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowStory,
}

func init() {
	storyCmd.AddCommand(showStoryCmd)
}
