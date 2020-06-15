package commands

import (
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
	db, err := cfg.GetDatabase(args[0])
	d.CheckError(err)
	defer db.Close()

	util.MergeStory(db, args[1], args[2], args[3], user)
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
