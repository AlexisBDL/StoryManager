package commands

import (
	"fmt"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/outputpager"

	"github.com/spf13/cobra"
)

var showTasks bool

func runShowStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	option := ""
	if showTasks {
		option = ".value.Tasks"
	}
	db, ds, err := cfg.GetPath(ID + option)
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
	Short: "Show the story ID.",
	Long:  "Show log whit '#commitID",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowStory,
}

func init() {
	storyCmd.AddCommand(showStoryCmd)

	showStoryCmd.Flags().BoolVarP(&showTasks, "tasks", "t", false, "Display tasks of story ID")
	showStoryCmd.Flags().Lookup("tasks").NoOptDefVal = "true"
}
