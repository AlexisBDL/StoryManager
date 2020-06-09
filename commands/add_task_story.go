package commands

import (
	"fmt"
	"strconv"

	"github.com/AlexisBDL/StoryManager/spec"
	"github.com/AlexisBDL/StoryManager/util"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"

	"github.com/spf13/cobra"
)

func runAddTaskStory(cmd *cobra.Command, args []string) error {
	ID := args[0]
	goal := args[1]
	maker := args[2]

	db, ds, err := cfg.GetDataset(ID)
	d.PanicIfError(err)
	defer db.Close()

	// Test Open
	_, valState, _ := cfg.GetPath(ID + storyState)
	if valState == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Story %s not found in my Stories", ID))
	}
	state, err := strconv.Unquote(types.EncodedValue(valState))
	d.PanicIfError(err)
	if state == stateClose {
		fmt.Printf("The story %s is close, you con't modify it\n", ID)
		return nil
	}

	// Edit
	var (
		fieldsList []string
	)

	resolvedList := cfg.ResolvePathSpec(ID) + storyTasks
	absPathTask := util.ApplyStructEdits(db, newTask(goal, maker, util.ListLen(resolvedList)), nil, nil)

	fieldsList = append(fieldsList, "@#"+absPathTask.Hash.String())
	absPathList := util.ListAppend(resolvedList, fieldsList)

	absPath, err := spec.NewAbsolutePath("#" + absPathList.Hash.String() + ".value")
	d.CheckError(err)

	absPathTask = util.ApplyStructEdits(db, newTask(goal, maker, util.ListLen(resolvedList)), nil, nil)

	// Commit
	title := getTitle(ID)
	msg := "Add task " + goal + " on story ID " + ID
	valPath := absPath.Resolve(db)

	util.Commit(db, ds, valPath, ID, msg, user, title)

	return nil
}

var addTaskStoryCmd = &cobra.Command{
	Use:   "Tadd <ID> <goal> <maker>",
	Short: "Add a task in story.",
	Args:  cobra.ExactArgs(3),
	RunE:  runAddTaskStory,
}

func init() {
	storyCmd.AddCommand(addTaskStoryCmd)
}
