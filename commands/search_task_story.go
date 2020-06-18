package commands

import (
	"github.com/AlexisBDL/StoryManager/util"

	"github.com/spf13/cobra"
)

var (
	goalTask  string
	makerTask string
	stateTask string
)

func runSearchTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

	resolvedList := cfg.ResolvePathSpec(ID) + storyTasks

	if goalTask != "" {
		util.ListGetBy(resolvedList, "Goal", goalTask)
	}
	if makerTask != "" {
		util.ListGetBy(resolvedList, "Maker", makerTask)
	}
	if stateTask != "" {
		util.ListGetBy(resolvedList, "State", stateTask)
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
	searchTaskStoryCmd.Flags().StringVarP(&goalTask, "goal", "g", "",
		"Search by maker",
	)
	searchTaskStoryCmd.Flags().StringVarP(&stateTask, "state", "s", "",
		"Search by state",
	)
}
