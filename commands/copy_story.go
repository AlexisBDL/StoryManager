package commands

import (
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util/datetime"

	"github.com/spf13/cobra"
)

var duplicate bool

func runCopyStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	dest := args[1]

	if duplicate {
		_, valTitle, err := cfg.GetPath(ID + storyTitle)
		title, err := strconv.Unquote(types.EncodedValue(valTitle))
		d.PanicIfError(err)
		data := []byte(title[:4] + datetime.Now().String()[20:28] + randomString(10))
		ID = hash.New(data[:20]).String()
		util.SyncStory("", ID, dest, cfg, true)
		fmt.Println("Duplicate ID is : " + ID)
	} else {
		util.SyncStory("", ID, dest, cfg, true)
	}

	return nil
}

var copyStoryCmd = &cobra.Command{
	Use:   "copy <ID> <destination>",
	Short: "Move the story <ID> with the databases <destination>.",
	Args:  cobra.ExactArgs(2),
	RunE:  runCopyStory,
}

func init() {
	storyCmd.AddCommand(copyStoryCmd)

	copyStoryCmd.Flags().BoolVarP(&duplicate, "duplicate", "d", false,
		"Duplicate story in my stories (new branch)",
	)
	copyStoryCmd.Flags().Lookup("duplicate").NoOptDefVal = "true"
}
