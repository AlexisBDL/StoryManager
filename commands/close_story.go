package commands

import (
	"fmt"

	"github.com/AlexisBDL/StoryManager/util"

	"github.com/attic-labs/noms/go/d"
	"github.com/spf13/cobra"
)

func runCloseStory(cmd *cobra.Command, args []string) error {
	ID := args[0]

	db, ds, err := cfg.GetDataset(ID)
	d.PanicIfError(err)
	defer db.Close()

	// Test Open
	if isOpenStory(ID) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID)
		return nil
	}

	// Edit close
	resolved := cfg.ResolvePathSpec(ID) + commitStory
	fields := []string{"State", stateClose}
	absPath := util.StoryEdit(db, resolved, fields)

	// Commit
	title := getTitle(ID)
	msg := "Story " + title + " with ID " + ID + " was closed"
	util.Commit(db, ds, absPath, ID, msg, user, title)

	return nil
}

var closeStoryCmd = &cobra.Command{
	Use:   "close <ID>",
	Short: "Close a story",
	Args:  cobra.ExactArgs(1),
	RunE:  runCloseStory,
}

func init() {
	storyCmd.AddCommand(closeStoryCmd)
}
