package commands

import (
	"fmt"

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
		switch stateTask {
		case "tr":
			util.ListGetBy(resolvedList, "State", stateTr)
		case "tt":
			util.ListGetBy(resolvedList, "State", stateTt)
		case "ec":
			util.ListGetBy(resolvedList, "State", stateEc)
		default:
			fmt.Printf("Your state it's not recognize, choose ec or tt or tr")
			return nil
		}
	}

	return nil
}

var searchTaskStoryCmd = &cobra.Command{
	Use:   "Tsearch <ID> [flag] <value>",
	Short: "Search task by value in story ID.",
	Long:  "For state, use flag -s and choose your state : ec -> Encours, tt -> Test, tr -> Terminé",
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
