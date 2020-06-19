package commands

import (
	"fmt"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func runUpdateStory(cmd *cobra.Command, args []string) error {
	dbTarget := args[0]

	dbT, err := cfg.GetDatabase(dbTarget)
	d.CheckError(err)
	defer dbT.Close()

	dbL, err := cfg.GetDatabase("")
	d.CheckError(err)
	defer dbL.Close()

	var (
		lsT []string
		lsL []string
		lsU []string
		ID  string
	)

	dbT.Datasets().IterAll(func(k, v types.Value) {
		ID = fmt.Sprint(k)
		lsT = append(lsT, ID)
	})

	dbL.Datasets().IterAll(func(k, v types.Value) {
		ID = fmt.Sprint(k)
		lsL = append(lsL, ID)
	})

	for _, v := range lsT {
		if !Find(lsL, v) {
			lsU = append(lsU, v)
		}
	}

	for _, v := range lsU {
		util.SyncStory(dbTarget+"::"+v, v, "Stories", cfg, false)
		title := getTitle(v)
		fmt.Printf("Story %s %s imported", v, title)
	}

	return nil
}

var updateStoryCmd = &cobra.Command{
	Use:   "update <dbTarget>",
	Short: "Add stories not imported in my database",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateStory,
}

func init() {
	RootCmd.AddCommand(updateStoryCmd)
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
