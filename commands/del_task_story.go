package commands

import (
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

func runDelTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	IDT := args[1]

	IDX, err := strconv.Atoi(IDT)
	d.PanicIfError(err)

	db, ds, err := cfg.GetDataset(ID)
	d.PanicIfError(err)
	defer db.Close()

	// Test Open
	if isOpenStory(ID) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID)
		return nil
	}

	// Edit
	resolvedList := cfg.ResolvePathSpec(ID) + storyTasks
	absPathDelT := util.ListDel(db, resolvedList, IDX)

	absPath, err := spec.NewAbsolutePath("#" + absPathDelT.Hash.String() + ".value")
	d.CheckError(err)

	// non termine, update ID / Fonction utile?????

	// Commit
	title := getTitle(ID)
	msg := "Del task " + IDT + " on story ID " + ID
	valPath := absPath.Resolve(db)

	util.Commit(db, ds, valPath, ID, msg, user, title)

	return nil
}

var delTaskStoryCmd = &cobra.Command{
	Use:   "Tdel <ID> <IDTask>",
	Short: "Del a task IDTask in story ID.",
	Args:  cobra.ExactArgs(2),
	RunE:  runDelTaskStory,
}

func init() {
	storyCmd.AddCommand(delTaskStoryCmd)
}
