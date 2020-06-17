package commands

import (
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

var (
	makerTask = "m"
	stateTask string
)

func runSearchTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

	db, err := cfg.GetDatabase(ID)
	d.PanicIfError(err)
	defer db.Close()

	if makerTask != "" {
		util.ListGetBy("Stories::"+ID+storyTasks, "Maker", makerTask)
	}
	if stateTask != "" {
		util.ListGetBy("Stories::"+ID+storyTasks, "State", stateTask)
	}

	return nil
}

var searchTaskStoryCmd = &cobra.Command{
	Use:   "Tsearch <ID> [flag] <value>",
	Short: "Search task by value in story ID.",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchTaskStory,
}

func init() {
	storyCmd.AddCommand(searchTaskStoryCmd)

	searchTaskStoryCmd.Flags().StringVarP(&makerTask, "maker", "m", "",
		"Search by maker",
	)
	searchTaskStoryCmd.Flags().StringVarP(&stateTask, "state", "s", "",
		"Search by state",
	)
}
