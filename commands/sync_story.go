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

	if isOpenStory(ID) {
		fmt.Printf("The story %s is close, you can't modify it\n", ID)
		return nil
	}

	if isOpenStory(destStore + "::" + ID) {
		fmt.Printf("The story %s is close in destination store, you can't modify it\n", ID)
		return nil
	}

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

	if choose {
		util.MergeStory(tempDB, "source", "dest", "merge", user, "p")
	} else if leftChoice {
		util.MergeStory(tempDB, "source", "dest", "merge", user, "l")
	} else if rightChoice {
		util.MergeStory(tempDB, "source", "dest", "merge", user, "r")
	}

	// sync temp --> merge ==> MyStory --> ID
	util.SyncStory("temp::merge", ID, "Stories", cfg, false)

	// sync temp --> merge ==> DestStory --> ID
	util.SyncStory("temp::merge", ID, destStore, cfg, true)

	tempDB.Close()
	os.RemoveAll("temp")

	return nil
}

var syncStoryCmd = &cobra.Command{
	Use:   "sync <ID> <destination/DbName> [Flag]",
	Short: "Syncronize the story <ID> with the databases <destination/DbName>. Choose the way to resolve conflicts whith flag. ID in my DB is left, ID in destination DB is right",
	Args:  cobra.ExactArgs(2),
	RunE:  runSyncStory,
}

func init() {
	storyCmd.AddCommand(syncStoryCmd)

	syncStoryCmd.Flags().BoolVarP(&choose, "choose", "c", false, "Ask me to choose between values if a conflict append")
	syncStoryCmd.Flags().Lookup("choose").NoOptDefVal = "true"

	syncStoryCmd.Flags().BoolVarP(&leftChoice, "left", "l", false, "Display open stories")
	syncStoryCmd.Flags().Lookup("left").NoOptDefVal = "true"

	syncStoryCmd.Flags().BoolVarP(&rightChoice, "right", "r", false, "Display open stories")
	syncStoryCmd.Flags().Lookup("right").NoOptDefVal = "true"
}
