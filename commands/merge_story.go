package commands

import (
	"fmt"
	"regexp"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"

	"github.com/spf13/cobra"
)

var (
	datasetRe = regexp.MustCompile("^" + datas.DatasetRe.String() + "$")
)

// Fusion de branch aillant le même parent
func runMergeStory(cmd *cobra.Command, args []string) error {
	ID1 := args[1]
	ID2 := args[2]
	IDM := args[3]

	db, err := cfg.GetDatabase(ID1)
	d.CheckError(err)
	defer db.Close()

	if isOpenStory(ID1) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID1)
		return nil
	}

	util.MergeStory(db, ID1, ID2, IDM, user)
	return nil
}

var mergeStoryCmd = &cobra.Command{
	Use:   "merge <ID1> <ID2> <merged>",
	Short: "Merge two stories that have similar ref.",
	Args:  cobra.ExactArgs(4),
	RunE:  runMergeStory,
}

func init() {
	storyCmd.AddCommand(mergeStoryCmd)
}
