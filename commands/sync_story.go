package commands

import (
	"fmt"
	"os"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"

	"github.com/spf13/cobra"
)

func runSyncStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	destStore := args[1]

	// temp --> source
	util.SyncStory(ID, "source", "temp", cfg, false)

	// temp --> dest
	if util.SyncStory(destStore+"::"+ID, "dest", "temp", cfg, true) {
		os.RemoveAll("temp")
		d.CheckErrorNoUsage(fmt.Errorf("Stories are already sync"))
		return nil
	}

	// merge
	tempDB, err := cfg.GetDatabase("temp")
	d.CheckError(err)
	defer tempDB.Close()

	util.MergeStory(tempDB, "source", "dest", "merge", user)

	// sync temp --> merge ==> MyStory --> ID
	util.SyncStory("temp::merge", ID, "Stories", cfg, false)

	// sync temp --> merge ==> DestStory --> ID
	util.SyncStory("temp::merge", ID, destStore, cfg, true)

	os.RemoveAll("temp")

	return nil
}

var pushStoryCmd = &cobra.Command{
	Use:   "sync <ID> <destination>",
	Short: "Syncronize the story <ID> with the databases <destination>.",
	Args:  cobra.ExactArgs(2),
	RunE:  runSyncStory,
}

func init() {
	storyCmd.AddCommand(pushStoryCmd)

}
