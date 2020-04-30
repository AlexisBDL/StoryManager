package commands

import (
	"fmt"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/outputpager"

	"github.com/spf13/cobra"
)

func runShowStory(cmd *cobra.Command, args []string) error {
	title := args[0]
	cfg := config.NewResolver() //config default db "Stories"
	db, val, err := cfg.GetPath(title)
	d.CheckError(err)
	defer db.Close()

	if val == nil {
		fmt.Printf("Story %s not found in database\n", title)
		return nil
	}

	pgr := outputpager.Start()
	defer pgr.Stop()

	types.WriteEncodedValue(pgr.Writer, val)
	fmt.Fprintln(pgr.Writer)

	return nil
}

var showStoryCmd = &cobra.Command{
	Use:   "show <title>",
	Short: "show a story.",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowStory,
}

func init() {
	storyCmd.AddCommand(showStoryCmd)
}
